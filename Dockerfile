FROM alpine/openssl AS generate-rsa-keys
WORKDIR /keys
RUN openssl genrsa -out keypair.pem 2048
RUN openssl rsa -in keypair.pem -pubout -out public.pem
RUN openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in keypair.pem -out private.pem

FROM maven:3.9.4-amazoncorretto-21-debian AS build
WORKDIR /build
RUN mkdir -p /build/.m2/repository
COPY pom.xml /build
RUN mvn -U -Dmaven.compiler.debug=true -Dmaven.compiler.debuglevel=lines,vars,source clean dependency:go-offline
COPY . /build
COPY --from=generate-rsa-keys /keys/*.pem /build/src/main/resources/certs/
RUN mvn -Dmaven.compiler.debug=true -Dmaven.compiler.debuglevel=lines,vars,source -DskipTests package
RUN ls /build/src/main/resources/certs/

FROM openjdk:21-jdk-slim AS run
WORKDIR /app
COPY --from=build /build/target/*.jar /app/project.jar
ENTRYPOINT java -agentlib:jdwp=transport=dt_socket,address=*:8081,server=y,suspend=n -jar project.jar -file log.txt
