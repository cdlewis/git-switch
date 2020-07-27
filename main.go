package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/manifoldco/promptui"
	"strings"
	"os/exec"
	"os"
)

func main() {
	gitDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	r, err := git.PlainOpen(gitDirectory)

	if err != nil {
		fmt.Println("open error!")
		panic(err)
	}

	branchesIter, err := r.Branches()
	defer branchesIter.Close()
	if err != nil {
		fmt.Println("branches error!")
		panic(err)
	}

	branches := []string{}
	err = branchesIter.ForEach(func (branch *plumbing.Reference) error {
		branches = append(branches, (*branch).Name().String()[11:])
		return nil
	})

	if err != nil {
		panic(err)
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
		Details: "{{ . }}",
	}

	prompt := promptui.Select{
		Label: "Branches",
		Items: branches,
		Size: 10,
		Searcher: searcher,
		HideSelected: true,
		Templates: templates,
	}

	_, branch, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// go-git's checkout functionality doesn't work properly
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = gitDirectory
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
