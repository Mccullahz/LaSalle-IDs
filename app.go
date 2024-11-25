package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// file dialog and pathing
func (a *App) SelectFile(dialogTitle string) string {
	selectedFile, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: dialogTitle,
	})
	if err != nil {
		fmt.Println("Error selecting file:", err)
		return ""
	}
	return selectedFile
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s :)", name)
}

func (a *App) GenerateIDCards(files map[string]string) string {
	// Currently only printing file paths, will add logic later
	fmt.Println("Data File:", files["dataFile"])
	fmt.Println("Student Codes File:", files["studentCodesFile"])
	fmt.Println("Staff Codes File:", files["staffCodesFile"])
	fmt.Println("ID Template File:", files["idTemplateFile"])

	return "ID cards generated successfully!"
}
