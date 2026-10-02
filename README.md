# OG-template - Datasektionen Go Starter Template

En ren, funktionell och produktionsredo mall för att snabbt bygga nya webbsystem inom **Konglig Datasektionen (KTH)**.

Mallen är inspirerad av [GOrdian](https://github.com/datasektionen/GOrdian) och innehåller allt som behövs för att komma igång på några minuter.

---

## Funktioner

- **Backend:** Go (standard `net/http` och `html/template` utan onödiga ramverk)
- **Autentisering & Accesser:** Delat paket [`GO-auth-yourself`](https://github.com/datasektionen/GO-auth-yourself) (`github.com/datasektionen/GO-auth-yourself`) för OIDC, JWT-sessioner och behörighetskontroll (`RequirePermissions`)
- **SSO / OIDC Inloggning:** Fullt integrerad med Datasektionens SSO (`sso.datasektionen.se` i prod, `nyckeln` i dev)
- **Design & Layout:**
  - **Methone Topbar:** Datasektionens standardmeny med dynamisk inloggningsknapp ("Logga in" / "Logga ut") och "Hem"
  - **Aurora CSS:** Datasektionens officiella designsystem
- **Databas:** PostgreSQL-stöd med automatisk inläsning av `schema.sql` (valfritt – OG-template startar även utan DB)
- **Lokal utveckling:** `docker compose` med `nyckeln-under-dorrmattan` (mock SSO med färdiga testanvändare och accesser)
- **Produktion:** Färdig `Dockerfile` (multi-stage Alpine med non-root user), `job.nomad.hcl` för Datasektionens Nomad-kluster och GitHub Actions CI/CD workflow
- **Enkelhet:** Minimal boilerplate med standardiserade Go- och Docker-kommandon

---

## Projektstruktur

```
OG-template/
├── main.go                       # Startpunkt och HTTP-server
├── schema.sql                    # Databasschema (körs vid start om DB är kopplad)
├── go.mod                        # Go moduler och beroenden
├── compose.yaml                  # Docker Compose med app, db och mock-SSO (nyckeln)
├── Dockerfile                    # Produktions-container
├── Dockerfile.dev                # Utvecklings-container
├── job.nomad.hcl                 # Nomad deploy-specifikation för produktion
├── .env.example                  # Exempel på miljövariabler
│
├── .github/
│   └── workflows/
│       └── deploy.yaml           # CI/CD: tester & bygg/push till ghcr.io
│
├── internal/
│   ├── config/
│   │   └── env.go                # Miljövariabelhantering och standardvärden
│   ├── database/
│   │   └── postgres.go           # PostgreSQL-anslutning och schema-inläsning
│   └── web/
│       ├── web.go                # Server struct, NewHandler() och route-registrering
│       ├── index.go              # Publika vyer (landningssida)
│       ├── admin.go              # Skyddade adminvyer och formulärhantering
│       ├── api.go                # JSON API endpoints (/api/health, /api/user)
│       ├── static/
│       │   └── main.css          # Egna stilmallar
│       └── templates/
│           ├── head.gohtml       # Methone & Aurora konfiguration
│           ├── footer.gohtml     # Datasektionen footer
│           ├── index.gohtml      # Landningssida
│           └── admin.gohtml      # Exempel på skyddad vy (kräver admin)
```

---

## Snabbstart (Lokal utveckling)

### Alternativ 1: Docker Compose med Live Reload (Rekommenderat)

Kör hela stacken inklusive mock-SSO, PostgreSQL och automatisk filsynkning:

```bash
docker compose watch
# eller: docker compose up --build
```

1. Öppna webbläsaren på [http://localhost:3000](http://localhost:3000).
2. Klicka på **"Logga in"** i toppmenyn (Methone).
3. Du omdirigeras till den lokala inloggningstjänsten (_Nyckeln_).
4. Välj eller skriv in en testanvändare:

| KTH-ID        | Namn          | Accesser   | Beskrivning / Syfte                                              |
| ------------- | ------------- | ---------- | ---------------------------------------------------------------- |
| **`turetek`** | Ture Teknolog | `admin`    | Full administratör. Har access till skyddade sidor som `/admin`. |
| **`diadat`**  | Diana Datalog | `view-all` | Användare med läsaccess(`view-all`).                             |
| **`marmed`**  | Martin Median | _(inga)_   | Standardinloggad användare utan särskilda accesser.              |

---

### Anpassa testanvändare och accesser lokalt

För att testa egna accesser (t.ex. `editor`, `kassor`, `bokningsansvarig`) kan du redigera `nyckeln.yaml`-sektionen i [compose.yaml](compose.yaml):

```yaml
hive:
  groups:
    - name: Redaktörer
      id: editors-group
      domain: example.com
      members:
        - turetek
      permissions:
        - id: editor # Läggs till i Tures SSO ID-token vid inloggning
```

---

### Alternativ 2: Köra direkt med `go run`

1. Starta stödtjänster (databas och mock-SSO):
   ```bash
   docker compose up db nyckeln -d
   ```
2. Kopiera miljövariabler och starta Go-servern:
   ```bash
   cp .env.example .env
   go run .
   ```

---

### Köra tester

Kör Go:s inbyggda testverktyg:

```bash
go test -v ./...
```

---

## Miljövariabler

| Variabel             | Beskrivning                                     | Standardvärde i dev                                                       |
| -------------------- | ----------------------------------------------- | ------------------------------------------------------------------------- |
| `SERVER_PORT`        | Port som servern lyssnar på                     | `3000`                                                                    |
| `SERVER_URL`         | Applikationens bas-URL                          | `http://localhost:3000`                                                   |
| `DB_CONN`            | PostgreSQL connection string (valfri)           | `postgres://postgres:postgres@localhost:5432/og_template?sslmode=disable` |
| `OIDC_PROVIDER`      | URL till OIDC Identity Provider                 | `http://localhost:7003` (`https://sso.datasektionen.se/op` i prod)        |
| `OIDC_CLIENT_ID`     | OAuth2 Client ID                                | `client-id`                                                               |
| `OIDC_CLIENT_SECRET` | OAuth2 Client Secret                            | `client-secret`                                                           |
| `OIDC_REDIRECT_URL`  | Callback-URL för OIDC                           | `http://localhost:3000/auth/callback`                                     |
| `APP_SECRET_KEY`     | Hemlig nyckel för att signera JWT-sessionskakor | Minst 32 tecken, appen startar inte annars (säker slumpsträng i prod)     |

---

## Hur du anpassar mallen för ditt nya system

Det krävs bara några enkla steg för att komma igång:

1. **Be D-Sys skapa ett nytt repo:**
   Kontakta D-Sys för att få ett nytt repository skapat under [github.com/datasektionen](https://github.com/datasektionen). Klona sedan det nya repot och kopiera in filerna från denna mall.

2. **Uppdatera modulnamn i [go.mod](go.mod):**

   ```bash
   go mod edit -module github.com/datasektionen/ditt-nya-system
   ```

   Ersätt alla imports i källkoden (`main.go`, `internal/web/web.go`) från `github.com/datasektionen/OG-template` till `github.com/datasektionen/ditt-nya-system` och städa beroenden med:

   ```bash
   go mod tidy
   ```

3. **Ändra systemnamn och färg i [head.gohtml](internal/web/templates/head.gohtml):**

   ```javascript
   window.methone_conf = {
       system_name: "Ditt System",
       color_scheme: "cerise", // "cerise", "blue_grey", "dark_blue", etc.
       ...
   };
   ```

4. **Databas & schema.sql (valfritt):**
   - **Om du vill använda databas:** Skriv dina tabeller i [schema.sql](schema.sql). De skapas automatiskt vid uppstart om `DB_CONN` är satt.
   - **Om du INTE vill ha databas:**
     - Ta bort miljövariabeln `DB_CONN` i `.env.example`, `compose.yaml` och `job.nomad.hcl`.
     - Ta bort filen [schema.sql](schema.sql) och mappen `internal/database/`.
     - Ta bort `database.Init` och `db`-parametern i [main.go](main.go).
     - Ta bort `db`-containern i [compose.yaml](compose.yaml) och kör `go mod tidy`.

5. **Lägg till nya sidor och mallar:**
   - Skapa `.gohtml`-mallar under [internal/web/templates/](internal/web/templates/).
   - Registrera routes och hanterare i [internal/web/web.go](internal/web/web.go).

6. **Skydda routes med accesser:**

   ```go
   // Kräver att användaren är inloggad och har accessen "admin" från SSO
   mux.Handle("/admin", s.authRoute(s.adminPage, "admin"))

   // Öppen sida (valfri inloggning)
   mux.Handle("/about", s.authRoute(s.aboutPage))
   ```

---

## Driftsättning (Datasektionen Nomad)

1. **Be D-Sys om OIDC-klient:**
   Be D-Sys registrera en ny OIDC-klient med redirect-URL:
   `https://dittsystem.datasektionen.se/auth/callback`

2. **Be D-Sys sätta upp GitHub Actions & Nomad-deployment:**
   Be D-Sys lägga in repot i Nomad-klustret och konfigurera deployment-secrets (`oidc_client_secret`, `app_secret_key`, eventuellt databaslösenord).

3. **Uppdatera [job.nomad.hcl](job.nomad.hcl):**
   Byt ut alla förekomster av `OG-template` mot ditt systemnamn:
   - `job "<ditt-system>"` och `group "<ditt-system>"`
   - `service { name = "<ditt-system>" ... }`
   - `traefik.http.routers.<ditt-system>.rule=Host(\`<ditt-system>.datasektionen.se\`)`
   - `traefik.http.routers.<ditt-system>.tls.certresolver=default`
   - `task "<ditt-system>"`
   - `SERVER_URL=https://<ditt-system>.datasektionen.se`
   - `OIDC_REDIRECT_URL=https://<ditt-system>.datasektionen.se/auth/callback`
   - `OIDC_CLIENT_ID=<ditt-system>`
   - `nomadVar "nomad/jobs/<ditt-system>"`
   - `DB_CONN=postgres://<ditt-system>:{{ .db_password }}@postgres.dsekt.internal:5432/<ditt-system>` _(eller ta bort raden om du inte har DB)_
   - `image_tag = "ghcr.io/datasektionen/<ditt-system>:latest"`

4. **Pusha till main:**
   När koden pushas till `main`-branschen i GitHub bygger GitHub Actions automatiskt Docker-imagen och pushar till GitHub Packages (`ghcr.io`). Nomad sköter därefter driftsättningen.
