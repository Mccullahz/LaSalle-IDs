package main

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx             context.Context
	DataFile        string
	IDTemplateFile  string
	OutputDirectory string
}
type Record struct {
	StudentID  string
	ImageName  string
	FirstName  string
	LastName   string
	CardNumber string
	LunchID    string
	IsStaff    bool
	Room       string // i dont think we need this but i remember seeing it somewhere
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) SelectDataFile() {
	selectedFile, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Data.csv File",
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV Files (*.csv)", Pattern: "*.csv"},
		},
	})
	if err != nil {
		fmt.Println("Error selecting Data.csv file:", err)
		return
	}
	a.DataFile = selectedFile
	runtime.EventsEmit(a.ctx, "dataFileSelected", selectedFile)
}
func (a *App) SelectIDTemplateFile() {
	selectedFile, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select ID Template File",
		Filters: []runtime.FileFilter{
			{DisplayName: "SVG Files (*.svg)", Pattern: "*.svg"},
		},
	})
	if err != nil {
		fmt.Println("Error selecting ID Template file:", err)
		return
	}
	a.IDTemplateFile = selectedFile
	runtime.EventsEmit(a.ctx, "idTemplateFileSelected", selectedFile)
}

func (a *App) SelectOutputDirectory() {
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Output Directory",
	})
	if err != nil {
		fmt.Println("Error selecting output directory:", err)
		return
	}
	a.OutputDirectory = selectedDir
	runtime.EventsEmit(a.ctx, "outputDirectorySelected", selectedDir)
}

func (a *App) parseDataCSV(filePath string) ([]Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	recordsData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []Record

	headers := recordsData[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	for _, row := range recordsData[1:] {
		record := Record{
			StudentID:  getField(row, headerMap, "Student ID"),
			ImageName:  getField(row, headerMap, "Image Name"),
			FirstName:  getField(row, headerMap, "First Name"),
			LastName:   getField(row, headerMap, "Last Name"),
			CardNumber: getField(row, headerMap, "CardNumber"),
			LunchID:    getField(row, headerMap, "Lunch ID"),
			IsStaff:    false,
		}

		//determine if the record is for a staff member
		if record.StudentID == "missing-Student ID" || strings.TrimSpace(record.StudentID) == "" {
			record.StudentID = "STAFF"
			record.IsStaff = true
		}

		// if LunchID is missing, set to "000000"
		if record.LunchID == "" {
			record.LunchID = "000000"
		}

		record.Room = ""

		records = append(records, record)
	}

	return records, nil
}

func getField(record []string, headerMap map[string]int, fieldName string) string {
	if index, ok := headerMap[fieldName]; ok && index < len(record) {
		return strings.TrimSpace(record[index])
	}
	return ""
}

func generateBarcode(data string) (string, error) {
	// basically ripped this from the powershell script because i dont want to fuck with the zint process
	cmd := exec.Command("./Zint/zint.exe", "--notext", "--nobackground", "--filetype=svg", "--direct", "--scale=0.6", "--height=30", "-b", "20", "-d", data)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error generating barcode: %v, output: %s", err, string(output))
	}
	svgContent := string(output)
	rectIndex := strings.Index(svgContent, "<rect")
	if rectIndex != -1 {
		svgContent = svgContent[rectIndex:]
	}
	return svgContent, nil
}

func getImageBase64(imagePath string) (string, error) {
	data, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str, nil
}

func sanitizeFilename(name string) string {
	invalidChars := []string{":", "/", "\\", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}
	return name
}

func (a *App) GenerateIDCards() string {
	if a.DataFile == "" {
		return "Please select the Data.csv file before generating ID cards."
	}
	if a.IDTemplateFile == "" {
		return "Please select the ID Template file before generating ID cards."
	}
	if a.OutputDirectory == "" {
		return "Please select an output directory before generating ID cards."
	}

	fmt.Println("Generating ID cards with the following files:")
	fmt.Println("Data File:", a.DataFile)
	fmt.Println("ID Template File:", a.IDTemplateFile)
	fmt.Println("Output Directory:", a.OutputDirectory)

	records, err := a.parseDataCSV(a.DataFile)
	if err != nil {
		log.Println("Error parsing Data.csv:", err)
		return "Error parsing Data.csv"
	}
	templateContent, err := ioutil.ReadFile(a.IDTemplateFile)
	if err != nil {
		log.Println("Error loading template:", err)
		return "Error loading template"
	}
	template := string(templateContent)

	outputDir := a.OutputDirectory
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Println("Error creating output directory:", err)
		return "Error creating output directory"
	}
	htmlContent := `<!DOCTYPE html><html><head><title>ID Cards</title><link rel="stylesheet" href="style.css" /></head><body>`

	for idx, record := range records {
		fmt.Printf("Processing %d/%d: %s %s\n", idx+1, len(records), record.FirstName, record.LastName)
		fmt.Print("First: ", record.FirstName, "Last: ", record.LastName, "\n")
		var barcodeSVG string
		if record.LunchID != "" && record.LunchID != "000000" {
			barcodeSVG, err = generateBarcode(record.LunchID)
			if err != nil {
				log.Printf("Error generating barcode for %s %s: %v\n", record.FirstName, record.LastName, err)
				barcodeSVG = ""
			}
		} else {
			log.Printf("No barcode for %s %s\n", record.FirstName, record.LastName)
			barcodeSVG = ""
		}
		imageBase64, err := getImageBase64("Images/" + record.ImageName)
		if err != nil {
			log.Printf("Error loading image for %s %s: %v\n", record.FirstName, record.LastName, err)
			imageBase64 = ""
		}
		output := template
		output = strings.ReplaceAll(output, "$base64", imageBase64)
		output = strings.ReplaceAll(output, "${First}", record.FirstName)
		output = strings.ReplaceAll(output, "${Last}", record.LastName)
		output = strings.ReplaceAll(output, "${id}", record.StudentID)
		output = strings.ReplaceAll(output, "${barcode}", barcodeSVG)

		safeFileName := sanitizeFilename(fmt.Sprintf("%s %s, %s.svg", record.StudentID, record.LastName, record.FirstName))

		outputFileName := filepath.Join(outputDir, safeFileName)

		err = os.WriteFile(outputFileName, []byte(output), 0644)
		if err != nil {
			log.Printf("Error writing SVG for %s %s: %v\n", record.FirstName, record.LastName, err)
			continue
		}

		htmlFilePath := filepath.Join(outputDir, "ID Cards.html")

		relativePath, err := filepath.Rel(filepath.Dir(htmlFilePath), outputFileName)
		if err != nil {
			relativePath = safeFileName
		}

		htmlContent += fmt.Sprintf(
			`<div class='page'><pre>%s</pre><iframe src="%s"></iframe><div class='pageNumber'>Page %d</div></div>`,
			record.CardNumber,
			relativePath,
			idx+1,
		)
	}

	htmlContent += "</body></html>"

	htmlFilePath := filepath.Join(outputDir, "ID Cards.html")
	err = os.WriteFile(htmlFilePath, []byte(htmlContent), 0644)
	if err != nil {
		log.Println("Error writing ID Cards.html:", err)
		return "Error writing ID Cards.html"
	}

	return "ID cards generated successfully!"
}
