#!/bin/bash

echo "Generating RSA keys..."

# Créer un dossier pour les clés si nécessaire
mkdir -p src/main/resources/certs

# Générer la paire de clés RSA
openssl genrsa -out src/main/resources/certs/keypair.pem 2048

# Extraire la clé publique
openssl rsa -in src/main/resources/certs/keypair.pem -pubout -out src/main/resources/certs/public.pem

# Créer une clé privée au format PKCS#8
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in src/main/resources/certs/keypair.pem -out src/main/resources/certs/private.pem

echo "RSA keys generated successfully in src/main/resources/certs/"