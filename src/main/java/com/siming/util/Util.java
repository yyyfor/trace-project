package com.siming.util;

import com.siming.user.CaEnrollment;
import com.siming.user.UserContext;
import org.apache.commons.io.IOUtils;
import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.openssl.PEMKeyPair;
import org.bouncycastle.openssl.PEMParser;
import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemReader;
import org.hyperledger.fabric.sdk.security.CryptoPrimitives;
import org.hyperledger.fabric.sdk.security.CryptoSuite;

import javax.xml.bind.DatatypeConverter;
import java.io.*;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.logging.Level;
import java.util.logging.Logger;

public class Util {

    private static KeyFactory kf;

    public static CaEnrollment getEnrollment(String keyFilePath, String keyFileName, String certFolderPath,
                                             String certFileName) throws IOException, NoSuchAlgorithmException, InvalidKeySpecException {
        PrivateKey key = null;
        String certificate = null;
        InputStream isKey = null;
        BufferedReader brKey = null;

        try {
            isKey = new FileInputStream(keyFilePath + File.separator + keyFileName);
            byte[] bytes = IOUtils.toByteArray(isKey);
            PemReader pr = new PemReader(new StringReader(new String(bytes)));
            PemObject po = pr.readPemObject();
            PEMParser pem = new PEMParser(new StringReader(new String(bytes)));

            if (po.getType().equals("PRIVATE KEY")) {
                key = new JcaPEMKeyConverter().getPrivateKey((PrivateKeyInfo) pem.readObject());
            } else {
                PEMKeyPair kp = (PEMKeyPair) pem.readObject();
                key = new JcaPEMKeyConverter().getPrivateKey(kp.getPrivateKeyInfo());
            }
            certificate = new String(Files.readAllBytes(Paths.get(certFolderPath, certFileName)));
        } finally {
            isKey.close();
//            brKey.close();
        }

        CaEnrollment enrollment = new CaEnrollment(key, certificate);
        return enrollment;
    }

    public static void cleanUp() {
        String directoryPath = "users";
        File directory = new File(directoryPath);
        deleteDirectory(directory);

    }

    public static boolean deleteDirectory(File dir) {
        if(dir.isDirectory()) {
            File[] children = dir.listFiles();
            for (int i = 0; i < children.length; i++) {
                boolean success = deleteDirectory(children[i]);
                if (!success) {
                    return false;
                }
            }
        }

        Logger.getLogger(Util.class.getName()).log(Level.INFO, "Deleteing - " + dir.getName());
        return dir.delete();
    }

    public static UserContext readUserContext(String affiliation, String username) throws IOException, ClassNotFoundException {
        String filePath = "users/" + affiliation + "/" + username + ".ser";
        File file = new File(filePath);
        if(file.exists()) {
            FileInputStream fileInputStream = new FileInputStream(file);
            ObjectInputStream in = new ObjectInputStream(fileInputStream);

            UserContext userContext = (UserContext) in.readObject();
            in.close();
            fileInputStream.close();
            return userContext;

        }

        return null;
    }

    public static void writeUserContext(UserContext userContext) throws IOException {
        String directoryPath = "users/" + userContext.getAffiliation();
        String filePath = directoryPath + "/" + userContext.getName() + ".ser";
        File directory = new File(directoryPath);
        if(!directory.exists()) {
            directory.mkdirs();
        }

        FileOutputStream fileOutputStream = new FileOutputStream(filePath);
        ObjectOutputStream out = new ObjectOutputStream(fileOutputStream);
        out.writeObject(userContext);
        out.close();
        fileOutputStream.close();
    }
}


