package com.siming.client;

import org.hyperledger.fabric.protos.peer.Chaincode;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.exception.ChaincodeEndorsementPolicyParseException;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;

import java.io.File;
import java.io.IOException;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.logging.Level;
import java.util.logging.Logger;

import static java.nio.charset.StandardCharsets.UTF_8;

public class ChannelClient {

    private String name;
    private Channel channel;
    private FabricClient fabricClient;

    public ChannelClient(String name, Channel channel, FabricClient fabricClient) {
        this.name = name;
        this.channel = channel;
        this.fabricClient = fabricClient;
    }

    public String getName() {
        return name;
    }

    public Channel getChannel() {
        return channel;
    }

    public FabricClient getFabricClient() {
        return fabricClient;
    }


    public Collection<ProposalResponse> sendTransactionProposal(TransactionProposalRequest request) throws InvalidArgumentException, ProposalException {
        Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,
                "Sending transaction proposal on channel " + channel.getName());
        Collection<ProposalResponse> responses = channel.sendTransactionProposal(request, channel.getPeers());
        for(ProposalResponse response: responses) {
            String stringResponse = new String(response.getChaincodeActionResponsePayload());
            Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,
                    "Transaction proposal on channel " + channel.getName() + " " + response.getMessage() + " "
                            + response.getStatus() + " with transaction id:" + response.getTransactionID());
            Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,stringResponse);
        }
        CompletableFuture<BlockEvent.TransactionEvent> cf = channel.sendTransaction(responses);
        Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,cf.toString());
        return responses;
    }

    public Collection<ProposalResponse> instantiateChainCode(String chaincodeName, String version, String chaincodePath, String language,
                                                             String functionName, String [] functionArgs, String policyPath) throws InvalidArgumentException, IOException, ChaincodeEndorsementPolicyParseException, ProposalException {
        Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,
                "Instantiate proposal request " + chaincodeName + " on channel " + channel.getName()
                        + " with Fabric client " + fabricClient.getInstance().getUserContext().getMspId() + " "
                        + fabricClient.getInstance().getUserContext().getName());
        InstantiateProposalRequest instantiateProposalRequest = fabricClient.getInstance().newInstantiationProposalRequest();
        instantiateProposalRequest.setProposalWaitTime(180000);
        ChaincodeID.Builder chaincodeIdBuilder = ChaincodeID.newBuilder().setName(chaincodeName).setVersion(version).setPath(chaincodePath);
        ChaincodeID chaincodeID = chaincodeIdBuilder.build();
        Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,
                "Instantiating Chaincode ID " + chaincodeName + " on channel " + channel.getName());
        instantiateProposalRequest.setChaincodeID(chaincodeID);
        if (language.equals(TransactionRequest.Type.GO_LANG.toString()))
            instantiateProposalRequest.setChaincodeLanguage(TransactionRequest.Type.GO_LANG);
        else {
            instantiateProposalRequest.setChaincodeLanguage(TransactionRequest.Type.JAVA);
        }
        instantiateProposalRequest.setFcn(functionName);
        instantiateProposalRequest.setArgs(functionArgs);
        Map<String, byte[]> map = new HashMap<>();
        map.put("HyperLedgerFabric", "InstantiateProposalRequest:JavaSDK".getBytes(UTF_8));
        map.put("method", "InstantiateProposalRequest".getBytes(UTF_8));
        instantiateProposalRequest.setTransientMap(map);

        if(policyPath != null) {
            ChaincodeEndorsementPolicy chaincodeEndorsementPolicy = new ChaincodeEndorsementPolicy();
            chaincodeEndorsementPolicy.fromYamlFile(new File(policyPath));
            instantiateProposalRequest.setChaincodeEndorsementPolicy(chaincodeEndorsementPolicy);
        }

        Collection<ProposalResponse> responses = channel.sendInstantiationProposal(instantiateProposalRequest);
        CompletableFuture<BlockEvent.TransactionEvent> cf = channel.sendTransaction(responses);
        Logger.getLogger(ChannelClient.class.getName()).log(Level.INFO,
                "Chaincode " + chaincodeName + " on channel " + channel.getName() + " instantiation " + cf);
        return responses;
    }
}
