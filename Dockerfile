FROM maven:3.9.4-amazoncorretto-21-debian AS build
WORKDIR /build
RUN mkdir -p /build/.m2/repository
COPY pom.xml /build
RUN mvn -U -Dmaven.compiler.debug=true -Dmaven.compiler.debuglevel=lines,vars,source clean dependency:go-offline
COPY . /build
RUN mvn -Dmaven.compiler.debug=true -Dmaven.compiler.debuglevel=lines,vars,source -DskipTests package

FROM openjdk:21-jdk-slim AS run
WORKDIR /app
COPY --from=build /build/target/*.jar /app/project.jar
ENTRYPOINT java -agentlib:jdwp=transport=dt_socket,address=*:8081,server=y,suspend=n -jar project.jar -file log.txt
