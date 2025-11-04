\# go-gamelist-gpt 🎮🤖



Ferramenta em \*\*Go\*\* para:



\- Ler um arquivo `gamelist.xml` (formato usado por EmulationStation, Batocera etc.)

\- Coletar o conteúdo das tags `<desc>` de cada jogo

\- Enviar essas descrições para a \*\*API GPT (OpenAI)\*\*

\- Gerar um \*\*novo XML\*\* (`gamelist\_pt.xml`) com as descrições \*\*traduzidas\*\* para português do Brasil



É um projeto com foco didático, mostrando na prática:



\- Como \*\*ler e escrever XML\*\* em Go

\- Como integrar com a \*\*API da OpenAI\*\* usando a lib `go-openai`

\- Como montar um pipeline simples de “ler arquivo → processar com IA → salvar resultado”



---



\## 🧠 Objetivo do projeto



A ideia principal não é só “traduzir gamelist”, mas:



\- Demonstrar \*\*como consumir a API ChatGPT\*\* em Go

\- Mostrar manipulação de \*\*estruturas XML reais\*\*

\- Criar um exemplo que qualquer pessoa que mexe com ROMs/emuladores entende e consegue reutilizar



Na prática, isso vira um:



> “Tradutor automático de descrições de jogos usando GPT”



---



\## ⚙️ Tecnologias utilizadas



\- \*\*Go (Golang)\*\*

&nbsp; - `encoding/xml`

&nbsp; - `os`, `path/filepath`

&nbsp; - `strings`, `time`, `fmt`, `io`

\- \*\*\[go-openai](https://github.com/sashabaranov/go-openai)\*\*  

&nbsp; Cliente não-oficial para a API da OpenAI em Go

\- \*\*API GPT (OpenAI Chat Completions)\*\*  

&nbsp; Usada para traduzir o conteúdo das tags `<desc>`.



---



\## 🗂️ Estrutura do XML (`gamelist.xml`)



O projeto espera um XML no padrão:



```xml

<gameList>

&nbsp; <game>

&nbsp;   <path>./roms/game1.zip</path>

&nbsp;   <name>Game 1</name>

&nbsp;   <desc>Descrição original em outro idioma...</desc>

&nbsp;   <!-- outros campos opcionais -->

&nbsp; </game>

&nbsp; <game>

&nbsp;   ...

&nbsp; </game>

</gameList>



