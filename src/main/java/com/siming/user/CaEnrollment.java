package com.siming.user;

import org.hyperledger.fabric.sdk.Enrollment;

import java.io.Serializable;
import java.security.PrivateKey;

public class CaEnrollment implements Enrollment, Serializable {

    private static final long serialVersionUID = 1234567890123L;
    private PrivateKey key;
    private String cert;


    public CaEnrollment(PrivateKey key, String cert) {
        this.key = key;
        this.cert = cert;
    }

    @Override
    public PrivateKey getKey() {
        return key;
    }

    @Override
    public String getCert() {
        return cert;
    }
}
