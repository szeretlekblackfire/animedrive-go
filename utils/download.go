package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
)

func FetchAndDownload(url string, headers map[string]string, baseFolder, ep, animeName string) {
	selectedQuality, selectedSourceURL, err := FetchAndSelectQuality(url, headers)
	if err != nil {
		fmt.Printf("❌ Error selecting quality: %v\n", err)
		return
	}

	animeFolder := filepath.Join(baseFolder, animeName)
	if _, err := os.Stat(animeFolder); os.IsNotExist(err) {
		err = os.Mkdir(animeFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("❌ Error creating folder: %v\n", err)
			return
		}
	}

	filename := fmt.Sprintf("%s - Episode %s - %s.mp4", animeName, ep, selectedQuality)
	filePath := filepath.Join(animeFolder, filename)
	err = DownloadFile(filePath, selectedSourceURL, headers)
	if err != nil {
		fmt.Printf("❌ Error downloading file: %v\n", err)
	} else {
		fmt.Println("✅ Download completed!")
	}
}

func FetchAndSelectQuality(url string, headers map[string]string) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	htmlData := string(body)
	sources := ScrapeSources(htmlData)
	if len(sources) == 0 {
		return "", "", fmt.Errorf("no sources found")
	}

	// Prompt user to select quality
	prompt := promptui.Select{
		Label: "🎥 Select Quality",
		Items: sources,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Quality }}",
			Active:   "🔥 {{ .Quality | cyan }}",
			Inactive: "  {{ .Quality | faint }}",
			Selected: "✅ {{ .Quality | red | cyan }}",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("prompt failed: %v", err)
	}

	selectedSource := sources[i]
	fmt.Printf("📥 Selected quality: %s\n", selectedSource.Quality)

	return selectedSource.Quality, selectedSource.URL, nil
}

func DownloadFile(filepath string, url string, headers map[string]string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a progress bar
	bar := pb.Full.Start64(resp.ContentLength)
	bar.SetTemplateString(`{{ red "🚀 Downloading:" }} {{counters . }} {{bar . "[" "=" ">" "_" "]"}} {{percent . }} {{speed . }}`)
	barReader := bar.NewProxyReader(resp.Body)
	defer bar.Finish()

	// Write the body to file
	_, err = io.Copy(out, barReader)
	if err != nil {
		return err
	}

	return nil
}
