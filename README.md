# GO-EXPERT
Repositório com os projetos da pós graduação go-expert da Fullcycle

# Desafio client-server-api

Este projeto implementa uma aplicação **cliente-servidor** em Go.
O servidor expõe uma API REST para consulta de cotações e o cliente consome essa API, salvando o resultado em arquivo.

##  Executando o servidor

1. Acesse a pasta do servidor:
   ```bash
   cd client-server-api/server
   ```

2. Execute o comando:
   ```bash
   go run server.go
   ```

3. Abra uma **nova janela ou aba do terminal**.

4. Na nova janela/aba, teste a API executando:
   ```bash
   curl http://localhost:8080/cotacao
   ```

5. Verifique que foi criado o arquivo **app.db** na pasta do servidor.

6. Verifique o conteúdo do banco com um cliente SQLite3, por exemplo:
   - Plugin do VS Code: **SQLite Viewer**
   - Ou qualquer outro cliente SQLite3 de sua preferência

##  Executando o cliente

1. Acesse a pasta do cliente:
   ```bash
   cd client-server-api/client
   ```

2. Execute o comando:
   ```bash
   go run client.go
   ```

3. Verifique que foi criado o arquivo **cotacao.txt** contendo a cotação do dólar.

# Desafio multithreading

Este projeto é uma aplicação multithread em Go que realiza consultas simultâneas a duas APIs de CEP: ViaCep e BrasilApi.
A aplicação seleciona a resposta mais rápida e a exibe no console, descartando a resposta mais lenta.
Se nenhuma das APIs responder em até 1 segundo, uma mensagem de timeout será exibida no console.

##  Executando a aplicação

1. Acesse a pasta do desafio:
   ```bash
   cd multithreading
   ```

2. Execute o comando:
   ```bash
   go run server.go
   ```

3. Verifique no console a resposta mais rápida.


