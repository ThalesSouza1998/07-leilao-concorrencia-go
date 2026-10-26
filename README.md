# Sistema de leilões com fechamento automático

Projeto do desafio do curso Go Expert (Full Cycle).

## Como o fechamento funciona

Ao criar um leilão, o registro é salvo com status `0` (aberto) e uma goroutine fica responsável pelo prazo. O `POST /auction` responde na hora: a espera acontece em segundo plano, sem prender a requisição.

Quando o tempo definido em `AUCTION_DURATION` termina, a goroutine atualiza o leilão no MongoDB para o status `1` (fechado). A atualização considera apenas leilões ainda abertos, então executá-la duas vezes não causa problema. A mesma duração é usada na validação dos lances, em `internal/infra/database/bid/create_bid.go`, para que não seja possível dar lance depois do encerramento.

O código do fechamento fica em `internal/infra/database/auction/create_auction.go`.

## Rodando com Docker

```bash
docker compose up --build
```

A API fica disponível em `http://localhost:8080` e o MongoDB em `localhost:27017`. O compose espera o banco ficar saudável antes de iniciar a aplicação.

Criando um leilão:

```bash
curl -i -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Notebook Dell",
    "category": "Informática",
    "description": "Notebook usado, com 16GB de memória",
    "condition": 2
  }'
```

A resposta é `201 Created`. Para descobrir o id gerado, consulte os leilões abertos:

```bash
curl "http://localhost:8080/auction?status=0"
```

Depois de `AUCTION_DURATION` (30 segundos no `.env` de exemplo), o mesmo leilão aparece com `"status": 1`:

```bash
curl http://localhost:8080/auction/SEU_ID_AQUI
```

## Rodando localmente

Com o MongoDB em execução (o comando abaixo sobe apenas o banco pelo compose):

```bash
docker compose up -d mongodb

export MONGODB_URL="mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
export MONGODB_DB=auctions
export AUCTION_DURATION=30s

go run ./cmd/auction
```

O arquivo `cmd/auction/.env` é carregado quando existe. Variáveis já exportadas no shell têm prioridade sobre ele, então o bloco acima serve para apontar a aplicação para o banco local em vez do host `mongodb` usado dentro do compose.

## Variáveis de ambiente

| Variável | Padrão | Descrição |
| --- | --- | --- |
| `AUCTION_DURATION` | `5m` | Tempo de duração do leilão. Aceita o formato do `time.ParseDuration`: `30s`, `10m`, `2h`. |
| `MONGODB_URL` | — | String de conexão do MongoDB. |
| `MONGODB_DB` | — | Nome do banco de dados. |
| `BATCH_INSERT_INTERVAL` | `3m` | Intervalo máximo para gravar o lote de lances. |
| `MAX_BATCH_SIZE` | `5` | Quantidade de lances que fecha um lote. |
| `MONGO_INITDB_ROOT_USERNAME` | — | Usuário root criado pelo container do MongoDB. |
| `MONGO_INITDB_ROOT_PASSWORD` | — | Senha do usuário root do MongoDB. |

O `.env` em `cmd/auction` já vem com valores prontos para o Docker Compose. Para mudar o tempo do leilão, basta alterar `AUCTION_DURATION` e subir os containers novamente.

## Endpoints

| Método | Rota | Descrição |
| --- | --- | --- |
| `POST` | `/auction` | Cria um leilão. |
| `GET` | `/auction?status=0` | Lista leilões por status, categoria e nome do produto. |
| `GET` | `/auction/:auctionId` | Busca um leilão pelo id. |
| `GET` | `/auction/winner/:auctionId` | Busca o leilão e o lance vencedor. |
| `POST` | `/bid` | Registra um lance. |
| `GET` | `/bid/:auctionId` | Lista os lances de um leilão. |
| `GET` | `/user/:userId` | Busca um usuário pelo id. |

Status do leilão: `0` para aberto e `1` para fechado.

## Testes

Os testes unitários rodam sem dependências externas:

```bash
go test ./...
```

O teste de integração do fechamento precisa do MongoDB. Ele é ignorado quando o banco não está acessível:

```bash
docker compose up -d mongodb

export MONGODB_URL="mongodb://admin:admin@localhost:27017/?authSource=admin"
go test ./... -v
```

O teste `TestFechamentoAutomaticoDoLeilao` cria um leilão com `AUCTION_DURATION=2s`, confirma que ele nasce aberto e consulta o banco até o status virar fechado. O teste `TestAgendamentoDoFechamentoNaoBloqueia` garante que o agendamento não segura a criação do leilão, e `TestDuracaoDoLeilao` cobre a leitura da variável de ambiente.
