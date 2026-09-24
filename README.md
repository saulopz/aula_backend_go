# Back-end RESTful em Go: exemplos das aulas

Exemplos da disciplina **Sistemas Distribuídos e Mobile**. Cada pasta é um projeto Go independente, que acompanha uma aula da wiki. A explicação completa de cada exemplo, linha a linha, está na página da aula:

**Wiki da disciplina:** https://wiki.arisa.com.br/index.php?title=Sistemas_Distribuídos_e_Mobile

Os exemplos formam uma sequência: cada um parte do anterior e acrescenta um conceito novo. Se você está começando, siga a ordem.

## Os exemplos

| Pasta | Aula | O que o exemplo acrescenta | Banco? |
|---|---|---|---|
| `01_hello` | Exemplo Hello World | servidor HTTP com a biblioteca padrão, rota com método (`GET /hello`) | não |
| `02_apikey` | Exemplo usando API KEY | autenticação por chave no cabeçalho, comparação segura, CORS | não |
| `03_apikeyuser` | API KEY por usuário | chaves por usuário lidas de um JSON, armadilha do `map` | não |
| `04_variaveis` | Manipulando uma variável | `PUT` com corpo, validação, condição de corrida e `sync.Mutex` | não |
| `05_json_structs` | JSON e Structs | JSON, `struct` com métodos, receptor ponteiro, arquivos separados | não |
| `06_colecoes` | Coleções, POST e parâmetros no caminho | coleção e item, `POST` com `201`, `{id}` no caminho, `map` | não |
| `07_usersapi` | Persistência com SQL | PostgreSQL com `database/sql`, padrão DAO, interfaces | sim |
| `08_usersapi_gorm` | Persistência com GORM | o mesmo serviço com o ORM GORM, trocando só o DAO | sim |
| `09_usersapi_seguranca` | Segurança na API | middlewares: API KEY com hash no banco, limite de taxa, limites e timeouts | sim |

## Pré-requisitos

- **Go 1.22 ou mais novo.** Confira com `go version`.
- **PostgreSQL**, para os exemplos 07 a 09. A wiki tem uma página de instalação e configuração para Linux (Ubuntu e Fedora), macOS e Windows: *PostgreSQL: Instalação e Configuração*.
- Para testar, a extensão **REST Client** (`humao.rest-client`) no VSCode ou no VSCodium. Cada pasta tem um arquivo `requests.http` com os testes.

## Como executar um exemplo

Dentro da pasta do exemplo:

```bash
go mod tidy
go run .
```

O `go mod tidy` baixa as dependências (só é necessário na primeira vez). O servidor fica ouvindo em `http://localhost:8080` até você encerrá-lo com `Ctrl+C`. Com ele rodando, abra o `requests.http` e clique em **Send Request** acima de cada teste.

Os exemplos usam sempre a porta 8080, então **rode um de cada vez**. Se aparecer `address already in use`, há outro exemplo ainda rodando.

Sempre que alterar o código, pare o servidor e rode de novo: ele não recarrega sozinho.

## Exemplos com API KEY (02 e 03)

As chaves desses exemplos ficam em arquivos que **não vão para o Git** (veja o `.gitignore`), porque chave não se publica. Crie-os na pasta do exemplo antes de rodar:

**`02_apikey/apikey.txt`**, com a chave na primeira linha:

```
01234
```

**`03_apikeyuser/apikeys.json`**:

```json
{
    "user1": "apikey1",
    "user2": "apikey2",
    "user3": "apikey3"
}
```

## Exemplos com banco de dados (07, 08 e 09)

Os três usam o mesmo banco. Crie-o uma vez, com o **seu** usuário do PostgreSQL (não com o `postgres`):

```bash
createdb usersapi
psql -d usersapi -f 09_usersapi_seguranca/schema.sql
```

O `schema.sql` do exemplo 09 contém todas as tabelas (a `users` e a `api_clients`), e pode ser aplicado de novo sem problemas. Os exemplos 07 e 08 usam só a `users`.

A conexão é configurada pela variável de ambiente `DATABASE_URL`, e nunca escrita no código:

```bash
# Linux e macOS
export DATABASE_URL="postgres://usuario:senha@localhost:5432/usersapi?sslmode=disable"

# Windows (PowerShell)
$env:DATABASE_URL="postgres://usuario:senha@localhost:5432/usersapi?sslmode=disable"
```

A variável vale só para o terminal onde foi definida.

### O exemplo 09 precisa de uma chave

No exemplo de segurança, todas as rotas (menos `GET /health`) exigem uma chave, gerada pelo próprio programa:

```bash
cd 09_usersapi_seguranca
go run . nova-chave meu-teste
```

Em `nova-chave meu-teste`, a primeira palavra é sempre a mesma, e a segunda é só um nome para identificar o cliente. **Não há senha a digitar**: o programa sorteia a chave e a mostra **uma única vez**. O banco guarda apenas o *hash*. Copie a chave para a variável `@chave` do `requests.http`, e **não faça commit** do arquivo com a chave colada.

## Versões do Go e das dependências

As versões mais recentes de algumas dependências exigem um Go mais novo:

| Dependência | Exemplos | Versão mais recente exige | Com Go mais antigo, use |
|---|---|---|---|
| `github.com/jackc/pgx/v5` | 07, 08, 09 | Go 1.25 | `go get github.com/jackc/pgx/v5@v5.7.4` |
| `gorm.io/driver/postgres` | 08, 09 | Go 1.25 | `go get gorm.io/driver/postgres@v1.6.0` |
| `golang.org/x/time` | 09 | Go 1.26 | `go get golang.org/x/time@v0.10.0` |

Com essas versões fixas, os exemplos funcionam no Go 1.22. O código é o mesmo nas duas situações.

Se aparecer um erro mencionando **toolchain**, é esse o motivo: o Go tentou baixar sozinho um compilador mais novo para atender a uma dependência. Use as versões da última coluna.

Cuidado ao levar um projeto entre computadores: a linha `go` do `go.mod` registra a versão exigida. Um projeto criado com Go 1.26 pede o 1.26 em qualquer máquina.

## Licença e como citar

O código deste repositório é distribuído sob a **licença MIT** (veja o arquivo [`LICENSE`](LICENSE)). Você pode usar, copiar, modificar e distribuir os exemplos livremente, inclusive em outros cursos e projetos, desde que mantenha o aviso de direitos autorais e a licença nas cópias.

Se este material for útil em suas aulas, cursos, trabalhos ou publicações, **cite-o**, por favor. Os dados para citação estão no arquivo [`CITATION.cff`](CITATION.cff). No GitHub, o botão **"Cite this repository"**, na lateral da página, gera a referência pronta. Em formato ABNT:

> ZAMBIASI, Saulo Popov. **Back-end RESTful em Go**: exemplos das aulas de Sistemas Distribuídos e Mobile. 2026. Disponível em: https://wiki.arisa.com.br/index.php?title=Sistemas_Distribuídos_e_Mobile. Acesso em: dia mês ano.