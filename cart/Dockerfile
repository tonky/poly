FROM mcr.microsoft.com/dotnet/core/sdk:3.0-alpine AS build
WORKDIR /app

COPY cart.fsproj .

RUN dotnet restore

COPY *.fs ./
COPY Templates ./Templates
COPY static ./static

RUN dotnet build -c Release --no-restore
RUN dotnet publish --no-build --nologo -c Release


FROM mcr.microsoft.com/dotnet/core/aspnet:3.0-alpine AS runtime
WORKDIR /app
ENV Logging__LogLevel__Default Warning

COPY --from=build /app/bin/Release/netcoreapp3.0 ./
ENTRYPOINT ["dotnet", "cart.dll"]
