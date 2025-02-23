# CHANGELOG

## Como Executar o Teste `create_auction_test.go`

### Pré-requisitos
- Certifique-se de ter o Go instalado na sua máquina.
- Certifique-se de ter o MongoDB instalado e em execução, ou use um servidor mock para testes.
- Configure as variáveis de ambiente necessárias.

### Passos para Executar o Teste

1. **Navegue até o Diretório do Teste:**
   Abra seu terminal e navegue até o diretório que contém o arquivo `create_auction_test.go`.
   Exemplo:
   ```sh
   cd internal/infra/database/auction
    ``` 
   
2. **Execute o Teste:**
   ```sh
   go test -v -run TestCreateAuction
   ```

   A flag -v é usada para exibir a saída do teste no terminal. 
   A flag -run é usada para especificar a função de teste a ser executada.

3. **View the Test Output:**
   Após executar o teste, você verá a saída no terminal. A saída do teste indicará se o teste passou ou falhou

4. **Interpreting the Test Output:**
   - Se o teste passar, você verá uma mensagem indicando que o teste passou.
   - Se o teste falhar, você verá uma mensagem de erro indicando o motivo da falha.