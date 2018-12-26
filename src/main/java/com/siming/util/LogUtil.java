package com.siming.util;

import com.siming.network.CreateChannel;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.Peer;

import java.util.Collection;
import java.util.Iterator;
import java.util.logging.Level;
import java.util.logging.Logger;

public class LogUtil {

    public static void channelLog(Channel channel) {
        Logger.getLogger(CreateChannel.class.getName()).log(Level.INFO, "Channel created "+channel.getName());
        Collection peers = channel.getPeers();
        Iterator peerIter = peers.iterator();
        while(peerIter.hasNext()) {
            Peer pr = (Peer) peerIter.next();
            Logger.getLogger(CreateChannel.class.getName()).log(Level.INFO,pr.getName()+ " at " + pr.getUrl());
        }
    }
}
