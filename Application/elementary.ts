//首先,建立与Gateway的gRPC连接
const peerEndpoint = 'localhost:7051';

async function newGrpcConnection(): Promise<grpc.Client> {
    const tlsRootCert = await fs.readFile(tlsCertPath);
    const tlsCredentials = grpc.credentials.createSsl(tlsRootCert);
    return new grpc.Client(peerEndpoint, tlsCredentials, {
        'grpc.ssl_target_name_override': 'peer0.org1.example.com',
    });
}
//其次,创建Gateway连接
const client = await newGrpcConnection();

const gateway = connect({
    client,
    identity: await newIdentity(),
    signer: await newSigner(),
});

async function newIdentity(): Promise<Identity> {
    const credentials = await fs.readFile(certPath);
    return { mspId: 'Org1MSP', credentials };
}

async function newSigner(): Promise<Signer> {
    const privateKeyPem = await fs.readFile(keyPath);
    const privateKey = crypto.createPrivateKey(privateKeyPem);
    return signers.newPrivateKeySigner(privateKey);
}
//第三，访问要调用的智能合约
const network = gateway.getNetwork(channelName);
const contract = network.getContract(chaincodeName);
const contract = network.getContract(chaincodeName, smartContractName);
//////

//第五，调用事务函数以读取和写入资产
const resultBytes = await contract.evaluateTransaction('GetAllAssets');

const resultJson = utf8Decoder.decode(resultBytes);
const result = JSON.parse(resultJson);
console.log('*** Result:', result);
//创建新资产
const assetId = `asset${Date.now()}`;

await contract.submitTransaction(
    'CreateAsset',
    assetId,
    'yellow',
    '5',
    'Tom',
    '1300',
);
//更新资产
const commit = await contract.submitAsync('TransferAsset', {
    arguments: [assetId, 'Saptha'],
});
const oldOwner = utf8Decoder.decode(commit.getResult());

console.log(`*** Successfully submitted transaction to transfer ownership from ${oldOwner} to Saptha`);
console.log('*** Waiting for transaction commit');

const status = await commit.getStatus();
if (!status.successful) {
    throw new Error(`Transaction ${status.transactionId} failed to commit with status code ${status.code}`);
}

console.log('*** Transaction committed successfully');
//查询更新后的资产
const resultBytes = await contract.evaluateTransaction('ReadAsset', assetId);

const resultJson = utf8Decoder.decode(resultBytes);
const result = JSON.parse(resultJson);
console.log('*** Result:', result);