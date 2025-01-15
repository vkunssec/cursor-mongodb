# Sistema de Paginação MongoDB com Cursor

Este projeto demonstra uma implementação eficiente de paginação em MongoDB usando a técnica de cursor baseado em ObjectID. O sistema está disponível em duas versões: Go e Node.js.

## 📋 Visão Geral

O sistema implementa uma solução de paginação que supera as limitações da paginação tradicional baseada em `skip/limit`, oferecendo melhor performance para grandes conjuntos de dados.

### Principais Características

- Paginação baseada em cursor usando ObjectID
- Monitoramento de comandos MongoDB
- Configuração via variáveis de ambiente
- Suporte para logging detalhado
- Tratamento de erros robusto

## 🚀 Tecnologias

### Versão Node.js
- Node.js
- MongoDB Driver para Node.js
- dotenv

### Versão Go
- Go
- MongoDB Driver para Go
- godotenv

## 💻 Implementação

### Paginação com Cursor

A paginação baseada em cursor funciona da seguinte forma:

1. Na primeira requisição, retorna os primeiros N documentos ordenados por `_id`
2. Nas requisições subsequentes, usa o último `_id` como ponto de referência
3. Busca documentos com `_id` maior que o último recebido

Vantagens desta abordagem:
- Performance consistente independente do tamanho da coleção
- Não afetado por inserções/deleções durante a paginação
- Consumo de memória eficiente

### Estrutura do Projeto
```
├── golang/
│   └── main.go
├── nodejs/
│   └── index.js
├── .env.example
└── README.md
```

## ⚙️ Configuração

1. Clone o repositório
2. Copie `.env.example` para `.env`
3. Configure as variáveis de ambiente:

```
MONGODB_URI="mongodb://localhost:27017/?retryWrites=true&w=majority"
MONGODB_DATABASE="sample_mflix"
MONGODB_COLLECTION="comments"
```

### Para Node.js
```bash
npm install
node nodejs/index.js
```

### Para Go
```bash
go mod tidy
go run golang/main.go
```

## 🔍 Monitoramento

O sistema inclui monitoramento detalhado das operações MongoDB:

- Logging de comandos enviados
- Registro de respostas recebidas
- Formatação JSON para melhor legibilidade
- Filtro para comandos não essenciais (ping, endSessions)

## 📊 Exemplo de Uso

```javascript
// Primeira página (primeiros 10 documentos)
const firstPage = await paginateWithCursor();

// Segunda página (próximos 10 documentos)
const lastId = firstPage[firstPage.length - 1]._id;
const secondPage = await paginateWithCursor(lastId);
```

## 🔒 Boas Práticas Implementadas

1. **Gerenciamento de Conexões**
   - Conexões são adequadamente fechadas
   - Timeout configurado para evitar bloqueios

2. **Tratamento de Erros**
   - Erro handling robusto
   - Logging apropriado de erros

3. **Configuração**
   - Uso de variáveis de ambiente
   - Configuração centralizada

4. **Performance**
   - Índices apropriados em `_id`
   - Paginação eficiente
   - Controle de recursos

## ⚠️ Considerações Importantes

- Certifique-se de ter índices adequados em `_id`
- Ajuste o tamanho do limite conforme necessidade
- Monitore o uso de memória em grandes conjuntos de dados
- Mantenha consistência na ordenação

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

