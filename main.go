package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var allowedFlags = []string{"--mark-done"}

func identifyHelp(command string) {
	switch command {
	case "add":
		HelpAdd()
	case "list":
		HelpList()
	case "delete":
		HelpDelete()
	case "done":
		HelpDone()
	default:
		fmt.Println(Colors.Format("Unknown command: "+command, Colors.Bold, Colors.Italic))
		fmt.Println("Available commands: add, list, done, delete")
	}
}

// create connection to Postgres database
func initDB(ctx context.Context) (*Store, error) {
	var connStr string

	// 1. Try OS env TODO_DB
	if v, ok := os.LookupEnv("TODO_DB"); ok && v != "" {
		connStr = v
		fmt.Println("Using TODO_DB from environment")
	}

	// 2. fallback to default
	if connStr == "" {
		connStr = "postgres://postgres:root@localhost:5432/postgres?sslmode=disable"
		fmt.Println("Using default connection string")
	}

	// 3. Try connecting
	store, err := NewStore(ctx, connStr)
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		// 3.1. Try running SQL file to create table
		fmt.Println("Attempting to initialize DB schema from table.sql…")
		if err := runSQLFile(ctx, connStr, "table.sql"); err == nil {
			fmt.Println("Schema initialized. Retrying store connection…")
			store, err = NewStore(ctx, connStr)
		}
	}

	// 4. Finally, if still error, prompt user
	if err != nil {
		fmt.Print("Enter DB connection string: ")
		fmt.Scanln(&connStr)
		store, err = NewStore(ctx, connStr)
	}

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return store, err
}

func checkMarkDoneFlag(input string) string {
	target := "--mark-done"

	// Exact match
	if input == target {
		return ""
	}

	// Simple Levenshtein distance (only for short strings so it’s fast)
	distance := func(a, b string) int {
		la, lb := len(a), len(b)
		dp := make([][]int, la+1)
		for i := range dp {
			dp[i] = make([]int, lb+1)
		}
		for i := 0; i <= la; i++ {
			for j := 0; j <= lb; j++ {
				if i == 0 {
					dp[i][j] = j
				} else if j == 0 {
					dp[i][j] = i
				} else if a[i-1] == b[j-1] {
					dp[i][j] = dp[i-1][j-1]
				} else {
					dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
				}
			}
		}
		return dp[la][lb]
	}

	if distance(input, target) <= 3 {
		return fmt.Sprintf("Did you mean %s?", target)
	}

	return "No such option available"
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	}
	return c
}

func main() {

	ctx := context.Background()
	store, err := initDB(ctx)

	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println(Colors.Format("No command provided. Use 'todo help' for usage instructions.", Colors.Bold, Colors.Red))
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "version", "--version", "-v", "--v":
		fmt.Println(Colors.Format("Todo CLI v1.1.3", Colors.Bold, Colors.Green, Colors.Underline))
		return

	case "help":
		if len(os.Args) < 3 {
			Help()
		}
		identifyHelp(strings.ToLower(os.Args[2]))

	case "list":
		todo_items, err := store.ListTodos(ctx, nil)
		if err != nil {
			fmt.Println("Failed to list todos:", err)
			os.Exit(1)
		}
		Table(todo_items)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println(Colors.Format(`Usage: todo add "<task>"`, Colors.Bold, Colors.Red))
			fmt.Println("Type 'todo help add' for more information.")
			os.Exit(1)
		}
		task := os.Args[2]
		if err := store.AddTodo(ctx, task); err != nil {
			fmt.Println("Failed to add todo:", err)
			os.Exit(1)
		}
		fmt.Println(Colors.Format("Todo added successfully!", Colors.Bold, Colors.Green))

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println(Colors.Format(`Usage: todo delete "<uid>"`, Colors.Bold, Colors.Red))
			fmt.Println("Type 'todo help delete' for more information.")
			os.Exit(1)
		}
		uid := os.Args[2]
		if uid == "all" {
			rows, err := store.DeleteAllTodos(ctx)
			if err != nil {
				fmt.Println("Failed to delete all todos:", err)
				os.Exit(1)
			}
			fmt.Println(Colors.Format(fmt.Sprintf("All %d todos deleted successfully!", rows), Colors.Bold, Colors.Green))
			return
		}
		if err := store.DeleteTodo(ctx, uid); err != nil {
			fmt.Println("Failed to delete todo:", err)
			os.Exit(1)
		}
		fmt.Println(Colors.Format("Todo deleted successfully!", Colors.Bold, Colors.Green))

	case "--mark-done":
		uid := os.Args[2]

		if uid == "all" {
			count, err := store.MarkAllDone(ctx)
			if err != nil {
				fmt.Println("Error in marking all todos as done!\n", err)
				os.Exit(1)
			}
			fmt.Printf("Marked %d todos as done!\n", count)
			return
		}
		if err := store.MarkDone(ctx, uid); err != nil {
			fmt.Println("Error in marking the todo as done!\n", err)
			os.Exit(1)
		}
		fmt.Println(Colors.Format("Todo marked done!", Colors.Bold, Colors.Green))

	case "--undone":
		if len(os.Args) < 3 {
			fmt.Println(Colors.Format(`Usage: todo --undone "<uid>"`, Colors.Bold, Colors.Red))
		}
		uid := os.Args[2]
		if uid == "all" {
			rows, err := store.UnmarkAllDone(ctx)
			if err != nil {
				fmt.Println("Error in unmarking all todos as undone!\n", err)
				os.Exit(1)
			}
			fmt.Printf("Unmarked %d todos as undone!\n", rows)
			return
		}

		if err := store.UnmarkDone(ctx, uid); err != nil {
			fmt.Println("Error in unmarking the todo as undone!\n", err)
			os.Exit(1)
		}
		fmt.Println(Colors.Format("Todo marked undone!", Colors.Bold, Colors.Green))

	default:
		if msg := checkMarkDoneFlag(command); msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Println("No such command")
			Help()
		}
		os.Exit(1)
	}
}

func checkFlags() {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			matched := false
			for _, good := range allowedFlags {
				// prefix match
				if strings.HasPrefix(good, strings.TrimPrefix(arg, "-")) {
					fmt.Printf("Did you mean %s?\n", good)
					os.Exit(1)
				}
			}
			if !matched {
				// regex fuzzy: match letters order
				re := regexp.MustCompile(strings.Join(strings.Split(arg, ""), ".*"))
				for _, good := range allowedFlags {
					if re.MatchString(good) {
						fmt.Printf("Did you mean %s?\n", good)
						os.Exit(1)
					}
				}
				fmt.Println("No such options available")
				os.Exit(1)
			}
		}
	}
}

func runSQLFile(ctx context.Context, connStr, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	parts := strings.Split(string(data), ";")
	store, err := NewStore(ctx, connStr)
	if err != nil {
		return err
	}
	defer store.Close()
	for _, stmt := range parts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := store.pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("exec failed on `%s`: %w", stmt, err)
		}
	}
	return nil
}

// Testing
func Test() {
	// ── 1) Connect to Postgres ───────────────────────────────────────────
	connStr := "postgres://postgres:root@localhost:5432/postgres?sslmode=disable"

	if connStr == "" {
		log.Fatal("set TODO_DB environment variable to your connection string")
	}

	ctx := context.Background()
	store, err := NewStore(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
	}
	defer store.Close()
	fmt.Println("✅ Connected to DB")

	// ── 2) Add a few sample todos ────────────────────────────────────────
	samples := []string{"Buy milk", "Walk the dog", "Write Go CLI"}
	for _, task := range samples {
		if err := store.AddTodo(ctx, task); err != nil {
			log.Fatalf("AddTodo(%q) failed: %v\n", task, err)
		}
		fmt.Printf("→ Added: %q\n", task)
	}

	// ── 3) List all todos ─────────────────────────────────────────────────
	fmt.Println("\nAll todos:")
	all, err := store.ListTodos(ctx, nil)
	if err != nil {
		log.Fatalf("ListTodos: %v\n", err)
	}
	Table(all)

	// ── 4) Mark the second task done ─────────────────────────────────────
	if len(all) >= 2 {
		uid := all[1].UID
		if err := store.MarkDone(ctx, uid); err != nil {
			log.Fatalf("MarkDone(%s): %v\n", uid, err)
		}
		fmt.Printf("\n✅ Marked done: %s (%s)\n", all[1].Task, uid)
	}

	// ── 5) List only done tasks ──────────────────────────────────────────
	done := true
	fmt.Println("\nDone todos:")
	doneList, err := store.ListTodos(ctx, &done)
	if err != nil {
		log.Fatalf("ListTodos(done=true): %v\n", err)
	}
	printTodos(doneList)

	// ── 6) Delete the first task ────────────────────────────────────────
	if len(all) >= 1 {
		uid := all[0].UID
		if err := store.DeleteTodo(ctx, uid); err != nil {
			log.Fatalf("DeleteTodo(%s): %v\n", uid, err)
		}
		fmt.Printf("\n🗑  Deleted UID: %s\n", uid)
	}

	// ── 7) Final list of all tasks ───────────────────────────────────────
	fmt.Println("\nFinal todos:")
	final, err := store.ListTodos(ctx, nil)
	if err != nil {
		log.Fatalf("ListTodos: %v\n", err)
	}
	printTodos(final)
}

func printTodos(todos []Todo) {
	for i, t := range todos {
		status := "❌"
		completed := "-"
		if t.Status {
			status = "✅"
			completed = t.CompletedAt.Format(time.RFC822)
		}
		fmt.Printf("%2d) [%s] %-15s added: %s  done: %s (%s)\n",
			i+1, t.UID, status,
			t.CreatedAt.Format("02-Jan-2006 3:04 PM"),
			completed,
			t.Task,
		)
	}
}
