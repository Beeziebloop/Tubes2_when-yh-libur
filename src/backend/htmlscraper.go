package main
import(
	"fmt"
	"io"
	"net/http" //ini dipake buat nge-fetch raw html dari url
	"strings"
	"time"
)

func fetchHTML(url string) (string, error){
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://"){
		return "", fmt.Errorf("url invalid!")
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp,err := client.Get(url)
	if err != nil{
		return "", fmt.Errorf("gagal fetch url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return "", fmt.Errorf("gagal baca response body: %w", err)
	}

	return string(body), nil
}

//basically dia yang handling ngeload html dari user input (baik dalam bentuk url atau mentahan)
func LoadHTML(input string) (*Node, error){
	input = strings.TrimSpace(input)
	if input == ""{
		return nil, fmt.Errorf("input kosong!")
	}

	var rawHTML string
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://"){
	   //fetch html dari url input
	   html, err := fetchHTML(input)
	   if err != nil{
		return nil, err
	   }
	   rawHTML = html
	}else{
		rawHTML = input
	}

	root := parseHTML(rawHTML)
	if root == nil{
		return nil, fmt.Errorf("gagal parsing html :(")
	}
	return root, nil
}