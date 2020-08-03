package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/manifoldco/promptui"
)

func main() {
	gitDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get current working directory")
		os.Exit(1)
	}

	r, err := git.PlainOpen(gitDirectory)
	if err != nil {
		fmt.Println("Unable to open repo. Git directory is possibly invalid.")
		os.Exit(1)
	}

	branchesIter, err := r.Branches()
	defer branchesIter.Close()
	if err != nil {
		fmt.Println("Unable to fetch branches for this project.")
		os.Exit(1)
	}

	branches := []string{}
	err = branchesIter.ForEach(func(branch *plumbing.Reference) error {
		branches = append(branches, (*branch).Name().String()[11:])
		return nil
	})
	if err != nil {
		fmt.Println("Unable to compile list of branches for this project.")
		os.Exit(1)
	}

	searcher := func(input string, index int) bool {
		branch := branches[index]
		return strings.Contains(branch, input)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "{{ . | cyan }}",
		Inactive: "{{ .Name | cyan }}",
		Selected: "{{ . | red | cyan }}",
		Details:  "{{ . }}",
	}

	prompt := promptui.Select{
		Label:        "Branches",
		Items:        branches,
		Size:         10,
		Searcher:     searcher,
		HideSelected: true,
		Templates:    templates,
	}

	_, branch, err := prompt.Run()
	if err != nil {
		fmt.Println("Unable to create prompt.")
		os.Exit(1)
	}

	// go-git's checkout functionality doesn't work properly
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = gitDirectory
	err = cmd.Run()
	if err != nil {
		fmt.Println("Unable to checkout selected branch.")
		os.Exit(1)
	}
}
