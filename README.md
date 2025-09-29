# i2l

`i2l`은 소스 코드를 분석하여 지식 그래프를 생성하고, 반대로 지식 그래프로부터 특정 프로그래밍 언어의 코드를 생성하는 도구입니다. 이 프로젝트는 Google의 Genkit 프레임워크를 사용하여 구현되었습니다.

## 주요 기능

-   **코드 분석 및 그래프 생성**: 소스 코드를 입력받아 코드 내의 엔티티(변수, 함수, 클래스 등)와 그들 간의 관계를 추출하여 지식 그래프(튜플 형태)를 생성합니다.
-   **그래프 기반 코드 생성**: 생성된 지식 그래프와 목표 프로그래밍 언어를 기반으로 새로운 소스 코드를 생성합니다.
-   **다중 AI 모델 지원**: Google Gemini와 같은 클라우드 기반 AI 모델뿐만 아니라, Ollama를 통해 로컬에서 실행되는 `gpt-oss`, `gemma`와 같은 모델도 지원합니다.

## 프로젝트 구조

-   `i2l.go`: `I2L` 핵심 로직이 담긴 파일입니다. Genkit 인스턴스와 AI 모델을 관리하며, `GenerateGraphFromCode`와 `GenerateCodeFromGraph`와 같은 주요 기능을 제공합니다.
-   `generate.go`: 코드로부터 그래프를 추출할 때 사용되는 프롬프트와 `Tuple` 구조체를 정의합니다.
-   `i2l/main.go`: `i2l` 라이브러리를 활용하는 예제 실행 파일입니다. 간단한 Go 코드를 그래프로 변환하고, 이 그래프를 기반으로 Java 코드를 생성하는 과정을 보여줍니다.
-   `models/`: Google AI, Ollama 등 다양한 AI 제공자로부터 특정 모델을 정의하고 가져오는 패키지입니다.
-   `go.mod`: 프로젝트의 모듈과 의존성을 정의합니다.

## 시작하기

### 사전 요구사항

-   Go (1.25.1 이상)
-   (선택 사항) 로컬 모델을 사용하려면 [Ollama](https://ollama.com/)를 설치해야 합니다.
-   (선택 사항) Google AI 모델을 사용하려면 `GEMINI_API_KEY`가 필요합니다.

### 설치

```bash
go build ./...
```

### 설정

Google AI 모델을 사용하려면, 환경 변수에 API 키를 설정해야 합니다.

```bash
export GEMINI_API_KEY="YOUR_API_KEY"
```

### CLI 도구 사용

`i2l/` 디렉토리에서 CLI 도구를 실행할 수 있습니다.

```bash
cd i2l
go run . -f <입력_파일> -l <대상_언어> -o <출력_파일> [-p <AI_공급자>]
```

**플래그 설명:**
- `-f`: 분석할 소스 코드 파일 경로
- `-l`: 생성할 대상 언어 (예: Java, Python, C#, etc.)
- `-o`: 생성된 코드를 저장할 출력 파일 경로
- `-p`: AI 공급자 선택 (`google` 또는 `ollama`, 기본값: `google`)

**사용 예제:**

```bash
# Google AI를 사용하여 Go 코드를 Java로 변환
go run . -f example.go -l Java -o result.java

# Ollama를 사용하여 Go 코드를 Python으로 변환
go run . -f example.go -l Python -o result.py -p ollama
```

실행 결과로 코드에서 추출된 그래프 튜플과, 이를 바탕으로 생성된 코드를 확인할 수 있습니다.

## 사용법

### 1. I2L 인스턴스 생성

Google AI 또는 Ollama를 사용하여 `I2L` 인스턴스를 초기화할 수 있습니다.

```go
// Google AI 사용
il, err := i2l.DefaultGoogleAIRAG(ctx)
if err != nil {
    panic(err)
}

// Ollama 사용
il, err := i2l.DefaultOllamaRAG(ctx)
if err != nil {
    panic(err)
}
```

### 2. 코드를 그래프로 변환

`GenerateGraphFromCode` 함수를 사용하여 소스 코드로부터 지식 그래프를 추출합니다.

```go
codeSnippet := `
package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}`

graph, err := il.GenerateGraphFromCode(ctx, codeSnippet)
if err != nil {
    panic(err)
}

for _, t := range graph {
    fmt.Println(t.String())
}
```

### 3. 그래프를 코드로 변환

`GenerateCodeFromGraph` 함수를 사용하여 지식 그래프로부터 원하는 언어의 코드를 생성합니다.

```go
codeResult, err := il.GenerateCodeFromGraph(ctx, "Java", graph)
if err != nil {
    panic(err)
}

fmt.Println(codeResult.Code)
```
