const { MongoClient, ObjectId } = require("mongodb");
const dotenv = require("dotenv");

// Carrega variáveis de ambiente
dotenv.config();

// Configuração do MongoDB
const config = {
	url: process.env.MONGODB_URI,
	appName: "sample_mflix",
	database: process.env.MONGODB_DATABASE,
	collection: process.env.MONGODB_COLLECTION,
};

// Função que implementa paginação usando cursor baseado em ObjectId
async function paginateWithCursor(lastId = null, limit = 10) {
	const client = await connect();

	try {
		const collection = client.db(config.database).collection(config.collection);

		// Cria filtro para buscar documentos após o último ID (se fornecido)
		const filter = lastId ? { _id: { $gt: new ObjectId(lastId) } } : {};

		// Configura opções de busca: limite e ordenação por _id
		const options = {
		limit,
		sort: { _id: 1 },
		};

		const results = await collection.find(filter, options).toArray();
		return results;
	} catch (error) {
		console.error("Erro na paginação:", error);
		throw error;
	}
}

// Função que estabelece conexão com MongoDB
async function connect() {
	const options = {
		appName: config.appName,
		monitorCommands: true, // Habilita monitoramento de comandos
	};

	const client = new MongoClient(config.url, options);

	// Configura monitor para logging de comandos MongoDB
	client.on("commandStarted", (event) => {
		if (!["endSessions", "ping"].includes(event.commandName)) {
			console.log(
				JSON.stringify(
					event,
					(_, value) => (typeof value === "bigint" ? value.toString() : value),
					2
				)
			);
		}
	});

	client.on("commandSucceeded", (event) => {
		if (!["endSessions", "ping"].includes(event.commandName)) {
			console.log(
				JSON.stringify(
					event,
					(_, value) => (typeof value === "bigint" ? value.toString() : value),
					2
				)
			);
		}
	});

	await client.connect();
	return client;
}

// Função principal que demonstra o uso da paginação
async function main() {
	try {
		// Primeira chamada: busca os primeiros 10 documentos
		const firstPage = await paginateWithCursor();
		console.log("Primeira página:", JSON.stringify(firstPage, null, 2));

		if (firstPage.length > 0) {
		// Segunda chamada: busca próximos 10 documentos após o último ID
		const lastId = firstPage[firstPage.length - 1]._id;
		const secondPage = await paginateWithCursor(lastId);
		console.log("Segunda página:", JSON.stringify(secondPage, null, 2));
		}
	} catch (error) {
		console.error("Erro:", error);
	} finally {
		// Fecha a conexão com MongoDB
		const client = await connect();
		await client.close();
		process.exit(0);
	}
}

// Executa o programa
main();
