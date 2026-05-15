# PassForge 🔐

> **Gerador de senhas criptográfico** — Go backend + interface web minimalista.

Aplicação web que usa `crypto/rand` do Go para gerar senhas com **aleatoriedade criptograficamente segura**. O servidor expõe um endpoint HTTP, e a interface (embutida no próprio binário) permite configurar comprimento, conjuntos de caracteres e exibe o nível de entropia da senha gerada.

## ✨ Funcionalidades

- 🔒 **Geração criptograficamente segura** — usa `crypto/rand` (sem `math/rand`)
- 📏 **Comprimento configurável** — de 4 a 64 caracteres
- 🔡 **4 conjuntos de caracteres** selecionáveis:
  - Maiúsculas (`A–Z`)
  - Minúsculas (`a–z`)
  - Números (`0–9`)
  - Símbolos (`!@#$%^&*` etc.)
- 📊 **Medidor de força** — calcula a entropia em bits e exibe nível: `WEAK → FAIR → GOOD → STRONG → FORTRESS`
- 📋 **Copiar com um clique** via Clipboard API
- 🎨 **UI dark/terminal** embutida diretamente no binário Go (sem arquivos estáticos externos)
- 🔀 **Fisher-Yates shuffle** criptográfico garante que nenhum conjunto seja ignorado

## 🚀 Como executar

### Pré-requisitos

- [Go](https://golang.org/dl/) 1.18 ou superior

### Rodando

```bash
# Clone o repositório
git clone https://github.com/buckgustavo/GENERATE-PASSWORD.git
cd GENERATE-PASSWORD

# Execute
go run main.go
```

Abra o navegador em **http://localhost:8080**

### Compilando

```bash
go build -o passforge main.go
./passforge
```

## 🌐 API

O servidor expõe um único endpoint REST:

### `POST /generate`

**Request body (JSON):**

```json
{
  "length": 16,
  "uppercase": true,
  "lowercase": true,
  "numbers": true,
  "symbols": false
}
```

**Response (JSON):**

```json
{
  "password": "aB3kRz9mTq2xLwN1",
  "entropy": 95,
  "poolSize": 62
}
```

| Campo | Descrição |
|---|---|
| `password` | Senha gerada |
| `entropy` | Entropia em bits \(= \text{length} \times \log_2(\text{poolSize})\) |
| `poolSize` | Tamanho do charset utilizado |
| `error` | Mensagem de erro (se houver) |

## 🛠️ Tecnologias

| Tecnologia | Uso |
|---|---|
| Go (`net/http`) | Servidor HTTP e lógica de geração |
| `crypto/rand` | Aleatoriedade criptograficamente segura |
| `encoding/json` | Serialização da API |
| HTML/CSS/JS (embutido) | Interface web servida pelo próprio Go |
| Fisher-Yates shuffle | Embaralhamento imparcial dos caracteres |

## 📁 Estrutura

```
GENERATE-PASSWORD/
└── main.go   # Servidor Go + HTML/CSS/JS da UI embutidos
```

## 🔢 Como a entropia é calculada

A entropia é calculada como:

\[ H = L \times \log_2(P) \]

Onde **L** é o comprimento da senha e **P** é o tamanho do pool de caracteres. Por exemplo, uma senha de 16 caracteres com maiúsculas + minúsculas + números tem entropia de ≈ 95 bits — classificada como **STRONG**.

| Bits | Nível |
|---|---|
| < 28 | WEAK |
| 28 – 49 | FAIR |
| 50 – 71 | GOOD |
| 72 – 99 | STRONG |
| ≥ 100 | FORTRESS |

---

Feito por [buckgustavo](https://github.com/buckgustavo)
