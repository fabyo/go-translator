package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// ATENÇÃO: coloque sua chave aqui SOMENTE PARA TESTE LOCAL.
const openAIKey = ""

// Estruturas básicas do gamelist.xml (EmulationStation)
type GameList struct {
	XMLName xml.Name `xml:"gameList"`
	Games   []Game   `xml:"game"`
}

type Game struct {
	Path        string `xml:"path"`
	Name        string `xml:"name"`
	Desc        string `xml:"desc"`
	Image       string `xml:"image,omitempty"`
	Rating      string `xml:"rating,omitempty"`
	Releasedate string `xml:"releasedate,omitempty"`
	Developer   string `xml:"developer,omitempty"`
	Publisher   string `xml:"publisher,omitempty"`
	Genre       string `xml:"genre,omitempty"`
	Players     string `xml:"players,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <pasta-ou-arquivo-gamelist.xml>")
		os.Exit(1)
	}

	if openAIKey == "" {
		fmt.Println("Configure a constante openAIKey com sua chave da OpenAI para testar.")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	// Se for pasta, monta caminho para gamelist.xml
	fi, err := os.Stat(inputPath)
	if err != nil {
		fmt.Println("Erro ao acessar caminho:", err)
		os.Exit(1)
	}

	xmlPath := inputPath
	if fi.IsDir() {
		xmlPath = filepath.Join(inputPath, "gamelist.xml")
	}

	fmt.Println("Lendo:", xmlPath)

	// Abre XML original
	f, err := os.Open(xmlPath)
	if err != nil {
		fmt.Println("Erro ao abrir XML:", err)
		os.Exit(1)
	}
	defer f.Close()

	var gl GameList
	if err := xml.NewDecoder(f).Decode(&gl); err != nil && err != io.EOF {
		fmt.Println("Erro ao decodificar XML:", err)
		os.Exit(1)
	}

	if len(gl.Games) == 0 {
		fmt.Println("Nenhum <game> encontrado no XML.")
		return
	}

	client := openai.NewClient(openAIKey)
	ctx := context.Background()

	fmt.Printf("Total de jogos: %d\n", len(gl.Games))

	for i := range gl.Games {
		orig := strings.TrimSpace(gl.Games[i].Desc)
		if orig == "" {
			continue // sem descrição, ignora
		}

		fmt.Printf("[%d/%d] Traduzindo '%s'...\n", i+1, len(gl.Games), gl.Games[i].Name)

		translated, err := translateDesc(ctx, client, orig)
		if err != nil {
			fmt.Println("  Erro na tradução, mantendo original:", err)
			continue
		}

		gl.Games[i].Desc = translated

		// só pra não bater que nem um doido na API
		time.Sleep(300 * time.Millisecond)
	}

	// Gera novo arquivo
	outPath := filepath.Join(filepath.Dir(xmlPath), "gamelist_pt.xml")
	if err := writeXML(outPath, &gl); err != nil {
		fmt.Println("Erro ao escrever XML traduzido:", err)
		os.Exit(1)
	}

	fmt.Println("XML traduzido gerado em:", outPath)
}

func translateDesc(ctx context.Context, client *openai.Client, text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return text, nil
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo, // pode trocar pelo modelo que você quiser
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Você é um tradutor de descrições de jogos. Traduza para português do Brasil mantendo o sentido e o estilo, sem acrescentar comentários.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		Temperature: 0.2,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("resposta vazia da API")
	}

	out := strings.TrimSpace(resp.Choices[0].Message.Content)
	return out, nil
}

func writeXML(path string, gl *GameList) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	enc := xml.NewEncoder(outFile)
	enc.Indent("", "  ")

	// cabeçalho XML
	if _, err := outFile.WriteString(xml.Header); err != nil {
		return err
	}

	return enc.Encode(gl)
}
