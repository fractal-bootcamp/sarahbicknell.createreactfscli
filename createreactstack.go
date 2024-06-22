package main 

import(
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

//variable to hold the available stack options
var stackOptions = []string{"React Vite + Express", "Remix"}

func main() {
	//creates the root command 

	var rootCmd = &cobra.Command{
		Use: "create-react-fs",
		//displays a description when it's run
		Short: "A generator for React fullstack projects",
		Long:  `This CLI tool helps you set up a React fullstack project with various options.`
		// run function defines what the command does when it's executed
		Run: func(cmd *cobra.Command, args []string) {
		// call selectStack function when create-react-fs is run 
			selectStack() 
		}
	}
	
	//common way to do error handling in go, this block calls Execute to run the root cmd and all its children
	//but the Execute method returns an error if something goes wrong, otherwise it returns nil 
	//so here we are checking that the returned value stored in err is not nil 
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		//exit the program with a code of 1 (non-zero exit codes by convention indicate an error)
		os.Exit(1)
	}
}

//function to let user select specific stack
func selectStack() {
	prompt := promptui.Select{
		Label: "Select Stack",
		Items: []string{"React Vite + Express", "Remix"}
	}

	index, result, err := prompt.Run() 
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	addDB, err := addDataBaseOptions() 
	if err != nil {
		fmt.Printf("Database option selection failed: %v\n", err)
		return
	}

	testing, err := addTestingOptions() 
	if err != nil {
		fmt.Printf("Testing selection failed: %v\n", err)
		return
	}

	fmt.Printf("You chose: %s\n", result)
	SetupStack(result)
	if addDB{
		setupDatabase() 
	}
	if testing{
		setupTesting() 
	}
}

//sets up the chosen stack
func SetupStack (stack string) {
	var cmd *exec.Cmd
	switch stack {
	case "React Vite + Express": 
		fmt.Println("Setting up React with Vite and an Express backend...")
		//shell commands for this
		cmd := exec.Command("npx", "create", "vite@latest", "--template", "react")
	case "Remix": 
		fmt.Println("Setting up a Remix project...")
		//setup cmnds
		cmd := exec.Command("npx", "create-remix@latest")
	default: 
		fmt.Println("Unknown stack option")
		return
	}

	//run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}
	fmt.Printf("Command output: %s\n", output)
	fmt.Println("Setup completed for", stack)
}

func addDataBaseOptions() (bool, error) {
	prompt := promptui.Select{
		Label: "Would you like to add a PostgresSQL database with Docker and Prisma?",
		Items: []string{"No", "Yes"},
	}
	
	index, _, err := prompt.Run() 
	if err != nil {
		return false, err
	}
	
	// return true if the user selected yes 
	return index == 1, nil
	
}

func setupDatabase() {
	fmt.Println("Starting PostgreSQL with Docker Compose...")
	//setup commands
	dockerComposeCmd := exec.Command("docker", "compose", "up")
	if err := dockerComposeCmd.Run(); err != nil {
		fmt.Printf("Failed to start Docker Compose: %v\n", err)
		return
	}

	fmt.Println("Installing Prisma as a dev dependency...")
	npmInstallPrismaCmd := exec.Command("npm", "install", "prisma", "--save-dev")
	if err := npmInstallPrismaCmd.Run(); err != nil { 
		fmt.Printf("Failed to install Prisma: %v\n", err)
		return
	}

	fmt.Println("PostgreSQL and Prisma setup completed!")
}

func addTestingOptions() (bool, error) {
	prompt := promptui.Select{
		Label: "Would you like to add testing with Vitest?",
		Items: []string{"No", "Yes"},
	}	

	index, _, err := prompt.Run() 
	if err != nil {
		return false, err
	}

	//return true if the user selected yes 
	return index == 1, nil 
}

func setupTesting() {
	fmt.Println("Installing vitest...")
	vitestCmd := exec.Command("npm", "install", "-D", "vitest")
	if err := vitestCmd.Run(); err != nil {
		fmt.Printf("Failed to install vitest: %v\n", err)
		return
	}
}

