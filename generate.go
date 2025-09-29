package i2l

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type Tuple struct {
	EntityX  string `json:"entity_x"`
	Relation string `json:"relation"`
	EntityY  string `json:"entity_y"`
}

func (t Tuple) String() string {
	return fmt.Sprintf("(%q, %q, %q)", t.EntityX, t.Relation, t.EntityY)
}

const generatePrompt = `# ROLE
You are an expert-level static code analysis engine. Your primary function is to act as a machine that constructs a code knowledge graph by analyzing source code.

# INSTRUCTION
Analyze the user-provided source code below. Identify all core code entities (e.g., modules, classes, functions, methods, variables, parameters) and the semantic relationships that connect them. Extract every meaningful relationship and represent them as a structured (Subject, Relationship, Object) tuple.

# RULES
1.  **Specificity**: The 'Relationship' part must be as specific and descriptive as possible, using verbs from the code's context. Examples include 'imports', 'defines', 'calls', 'inherits from', 'implements', 'instantiates', 'assigns to', 'returns', 'has parameter'. Avoid generic terms like "is related to".
2.  **Normalization**: Standardize entities to their most recognizable form (e.g., use fully qualified names or clear, context-based names like 'ClassName.methodName'). Resolve pronouns or contextual references where applicable.
3.  **Completeness**: Extract all possible valid relationships from the code snippet. A single line of code can contain multiple tuples.
4.  **Focus on Semantics**: Ignore non-semantic elements like comments. Analyze string literals only when they are an object in a clear relationship (e.g., a function printing or returning a specific string).
5.  **Format**: The final output must be a Python list of tuples. If no relationships are found, return an empty list [].

# EXAMPLES
-   **Text**:
    '''python
    import pandas as pd

    class DataProcessor:
        def __init__(self, file_path):
            self.df = pd.read_csv(file_path)

        def clean_data(self):
            return self.df.dropna()
    '''
-   **Output**:
    '''python
    [
        ("DataProcessor", "imports", "pandas as pd"),
        ("DataProcessor", "defines method", "__init__"),
        ("__init__", "has parameter", "file_path"),
        ("__init__", "calls", "pd.read_csv"),
        ("pd.read_csv", "is assigned to", "self.df"),
        ("DataProcessor", "defines method", "clean_data"),
        ("clean_data", "calls", "self.df.dropna"),
        ("clean_data", "returns", "self.df.dropna()")
    ]
    '''

-   **Text**:
    '''javascript
    import { apiFetch } from './utils.js';

    async function getUser(userId) {
        const userData = await apiFetch('/users/${userId}');
        console.log(userData.name);
    }
    '''
-   **Output**:
    '''python
    [
        ("getUser", "imports", "apiFetch from ./utils.js"),
        ("getUser", "is an async function", "True"),
        ("getUser", "has parameter", "userId"),
        ("getUser", "defines variable", "userData"),
        ("getUser", "calls", "apiFetch"),
        ("apiFetch", "is assigned to", "userData"),
        ("getUser", "calls", "console.log"),
        ("console.log", "accesses property", "userData.name")
    ]
    '''
-   **Text**: '// This is just a comment.'
-   **Output**: []

---

# TASK
Analyze the following source code and extract the relationships according to the rules and examples provided above.

`

func (i2l *I2L) GenerateGraphFromCode(ctx context.Context, code string) ([]Tuple, error) {
	response, _, err := genkit.GenerateData[[]Tuple](ctx, i2l.g, ai.WithModel(i2l.classificationModel), ai.WithPrompt(generatePrompt+code))
	if err != nil {
		return nil, fmt.Errorf("failed to generate graph from code: %w", err)
	}

	return *response, nil
}

type CodeGenerationResult struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

const codeGenPrompt = `# ROLE
You are an expert-level Code Synthesis Engine. Your primary function is to act as a machine that reconstructs source code from a semantic knowledge graph, represented as a list of relationship tuples.

# INSTRUCTION
Analyze the user-provided list of (Subject, Relationship, Object) tuples and a target programming language. Your task is to reconstruct the original source code from these relationships. The generated code must be syntactically correct, logically ordered, properly formatted, and idiomatic for the specified target language.

# RULES
1.  **Language Adherence**: Strictly generate code that conforms to the syntax, conventions, and best practices of the specified 'Target Language'. For example, use correct indentation for Python, curly braces and semicolons for JavaScript/Java, etc.
2.  **Code Formatting**: Apply proper code formatting and styling according to the target language's standards:
    -   Use consistent indentation (4 spaces for Python, 2 or 4 spaces for JavaScript/TypeScript, etc.)
    -   Apply proper line spacing between functions, classes, and logical blocks
    -   Follow naming conventions (camelCase, snake_case, PascalCase as appropriate)
    -   Use appropriate whitespace around operators and after commas
    -   Ensure proper bracket and parentheses alignment
    -   Add appropriate line breaks for readability
3.  **Logical Ordering**: The input tuples may not be in a sequential order. You must infer the correct sequence of code. For instance, 'imports' must be placed at the top of the file, class and function definitions must precede their use, and variables must be declared or defined before being referenced.
4.  **Relationship Interpretation**: Accurately translate the semantic relationships into concrete code constructs.
    -   ('A', 'imports', 'B') -> 'import B as A' or 'const A = require('B')'
    -   ('MyClass', 'defines method', 'myMethod') -> 'class MyClass { myMethod() {...} }'
    -   ('myMethod', 'has parameter', 'param') -> 'def myMethod(self, param):'
    -   ('myMethod', 'calls', 'otherFunc') -> The body of 'myMethod' should contain a call to 'otherFunc()'.
    -   ('value', 'is assigned to', 'variable') -> 'variable = value'
5.  **Completeness**: Attempt to incorporate all provided tuples into the final code. If a relationship is ambiguous or conflicts with another, make a logical assumption that best fits the overall context and structure.
6.  **Format**: The final output must be a single, clean, well-formatted code block containing only the generated source code for the specified language. Do not add any explanations or comments outside the code.

# EXAMPLES
-   **Input**:
    -   **Tuples**:
        '''python
        [
            ("DataProcessor", "imports", "pandas as pd"),
            ("DataProcessor", "defines method", "__init__"),
            ("__init__", "has parameter", "file_path"),
            ("__init__", "calls", "pd.read_csv"),
            ("pd.read_csv", "is assigned to", "self.df"),
            ("DataProcessor", "defines method", "clean_data"),
            ("clean_data", "calls", "self.df.dropna"),
            ("clean_data", "returns", "self.df.dropna()")
        ]
        '''
    -   **Target Language**: 'Python'
-   **Output**:
    '''python
    import pandas as pd

    class DataProcessor:
        def __init__(self, file_path):
            self.df = pd.read_csv(file_path)

        def clean_data(self):
            return self.df.dropna()
    '''

-   **Input**:
    -   **Tuples**:
        '''python
        [
            ("getUser", "imports", "apiFetch from ./utils.js"),
            ("getUser", "is an async function", "True"),
            ("getUser", "has parameter", "userId"),
            ("getUser", "defines variable", "userData"),
            ("getUser", "calls", "apiFetch"),
            ("apiFetch", "is assigned to", "userData"),
            ("getUser", "calls", "console.log"),
            ("console.log", "accesses property", "userData.name")
        ]
        '''
    -   **Target Language**: 'JavaScript'
-   **Output**:
    '''javascript
    import { apiFetch } from './utils.js';

    async function getUser(userId) {
        const userData = await apiFetch('/users/${userId}');
        console.log(userData.name);
    }
    '''

---

# TASK
Based on the provided list of relationship tuples and the target language, generate the corresponding source code.

`

func (i2l *I2L) GenerateCodeFromGraph(ctx context.Context, lang string, graph []Tuple) (CodeGenerationResult, error) {
	prompt := codeGenPrompt + "language: " + lang + "\nTuples:\n"
	for _, t := range graph {
		prompt += t.String() + "\n"
	}
	response, _, err := genkit.GenerateData[CodeGenerationResult](ctx, i2l.g, ai.WithModel(i2l.generativeModel), ai.WithPrompt(prompt))
	if err != nil {
		return CodeGenerationResult{}, fmt.Errorf("failed to generate code from graph: %w", err)
	}

	return *response, nil
}
