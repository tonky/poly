FROM mcr.microsoft.com/dotnet/core/sdk:3.0 AS build-env
WORKDIR /app

# Copy csproj and restore as distinct layers
# COPY *.csproj ./
COPY cart.fsproj ./
# COPY src/Migrations/Migrations.fsproj ./
# RUN dotnet restore Migrations.fsproj
RUN dotnet restore

# Copy everything else and build
# COPY src/cart/* ./
# RUN dotnet publish -c Release -o out
COPY . ./
RUN dotnet publish -c Release -o out

# Build runtime image
FROM mcr.microsoft.com/dotnet/core/runtime:3.0
# FROM alpine:latest
WORKDIR /app
COPY --from=build-env /app/out .
ENV Logging__LogLevel__Default Warning
ENTRYPOINT ["dotnet", "cart.dll"]
