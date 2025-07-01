package main

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Printf("%s%s todo %s", Colors.Format("Usage", Colors.Bold, Colors.Underline, Colors.Green), Colors.Format(": ", Colors.Bold, Colors.Green), Colors.Format("[command] [options]", Colors.Bold, Colors.Red))
	fmt.Println(Colors.Format("\n\nCommands:", Colors.Bold, Colors.Underline))
	fmt.Println("  add <task>           Add a new todo task")
	fmt.Println("  list                 List all todo tasks")
	fmt.Println("  done <uid>           Mark a task as done by its UID")
	fmt.Println("  delete <uid>         Delete a task by its UID")
	fmt.Println("  help                 Show this help message")
	fmt.Println(Colors.Format("\nOptions:", Colors.Bold, Colors.Underline))
	fmt.Println("  --done               List only done tasks")

	fmt.Println("")
	os.Exit(0)
}

func HelpAdd() {
	fmt.Println(Colors.Format("Usage: ", Colors.Bold, Colors.Green), Colors.Format("todo add [task]", Colors.Bold, Colors.Red))
	fmt.Println(Colors.Format("\nDescription:", Colors.Bold, Colors.Underline))
	fmt.Println("  Adds a new todo task.")
	// fmt.Println(Colors.Format("\nOptions:", Colors.Bold, Colors.Underline))
	// fmt.Println("  --priority <level>  Set the priority level (low, medium, high)")

	fmt.Println("")
	os.Exit(0)
}

func HelpList() {
	fmt.Println(Colors.Format("Usage: ", Colors.Bold, Colors.Green), Colors.Format("todo list [options]", Colors.Bold, Colors.Red))
	fmt.Println(Colors.Format("\nDescription:", Colors.Bold, Colors.Underline))
	fmt.Println("  Lists all todo tasks.")
	fmt.Println(Colors.Format("\nOptions:", Colors.Bold, Colors.Underline))
	fmt.Println("  --done          List only done tasks")

	fmt.Println("")
	os.Exit(0)
}

func HelpDone() {
	fmt.Println(Colors.Format("Usage: ", Colors.Bold, Colors.Green), Colors.Format("todo done <uid>", Colors.Bold, Colors.Red))
	fmt.Println(Colors.Format("\nDescription:", Colors.Bold, Colors.Underline))
	fmt.Println("  Marks a task as done by its UID.")
	fmt.Println(Colors.Format("\nOptions:", Colors.Bold, Colors.Underline))
	fmt.Println("  --force         Force mark as done without confirmation")

	fmt.Println("")
	os.Exit(0)
}

func HelpDelete() {
	fmt.Println(Colors.Format("Usage: ", Colors.Bold, Colors.Green), Colors.Format("todo delete <uid>", Colors.Bold, Colors.Red))
	fmt.Println(Colors.Format("\nDescription:", Colors.Bold, Colors.Underline))
	fmt.Println("  Deletes a task by its UID.")

	fmt.Println("")
	os.Exit(0)
}

// func main() {
// 	Help()
// 	fmt.Println("--------------------------------------------------------------------------")
// 	HelpAdd()
// 	fmt.Println("--------------------------------------------------------------------------")
// 	HelpList()
// 	fmt.Println("--------------------------------------------------------------------------")
// 	HelpDone()
// 	fmt.Println("--------------------------------------------------------------------------")
// }
