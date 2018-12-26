package com.siming;

import com.siming.config.Config;

import java.io.File;

public class Test {

    public static void main(String[] args) {
        File file = new File(Config.CHANNEL_CONFIG_PATH);
        File [] files = file.listFiles();
        for (File aa: files) {
            System.out.println(aa.getName());
        }

    }
}
