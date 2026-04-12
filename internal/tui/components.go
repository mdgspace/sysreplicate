package tui

import (
	"fmt"
	"os"
)

func PrintText(s []string) {
	for _, str := range s {
		fmt.Println(textStyle.Render(str))
	}
}

func PrintTitle(s []string) {
	for _, str := range s {
		fmt.Println(titleStyle.Render(str))
	}
}

func PrintError(s []string) {
	for _, str := range s {
		fmt.Fprintln(os.Stderr, errorStyle.Render(str))
	}
	os.Exit(1)
}

func RenderPrimaryBanner(str string) string {
	return primaryBanner.Render(str)
}

func RenderSecondaryBanner(str string) string {
	return secondaryBanner.Render(str)
}

func RenderFooter(str string) string {
	return footerStyle.Render(str)
}

func RenderListItem(index int, cursorIndex int, selectedIndex int, choice string) string {
	cursorStr := " "
	if cursorIndex == index {
		cursorStr = cursorStyle.Render(">")
	}

	checkedStr := " "
	if selectedIndex == index+1 {
		checkedStr = checkedStyle.Render("x")
	}

	line := fmt.Sprintf("%s [%s] ", cursorStr, checkedStr)
	return itemStyle.Render(line) + textStyle.Render(choice) + "\n"
}
