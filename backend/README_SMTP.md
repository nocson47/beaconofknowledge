SMTP setup and testing

This project supports sending password-reset emails either via a console logger (development) or through a real SMTP server.

Configuration

Add the following environment variables (or put them in `.env`) to enable SMTP:

- SMTP_HOST: host of SMTP server (e.g. smtp.gmail.com)
- SMTP_PORT: port (587 for STARTTLS, 465 for SMTPS)
- SMTP_USER: username for SMTP (optional; if empty authentication is skipped)
- SMTP_PASS: password for SMTP user (optional)
- SMTP_FROM: from address to use in emails (e.g. "YourApp <no-reply@example.com>")

How it works

- If SMTP_HOST and SMTP_PORT are set, the backend will instantiate the SMTP sender and use it to send reset emails.
- If not set, the backend uses `ConsoleEmailSender` which prints reset links to the application log (look for `[ConsoleEmail]` entries).

Testing locally

1) Using console logger (default/dev):
   - Start backend normally (e.g. `go run ./cmd` or `docker compose up backend`).
   - Trigger a forgot password from the frontend.
   - Watch backend logs for a line containing `[ConsoleEmail]` and `ResetURL=`. Copy that URL into your browser to open the Reset page.

2) Using a real SMTP server:
   - Populate `.env` with the SMTP_* variables above.
   - Restart backend.
   - Trigger forgot password from frontend; an email should be delivered to the provided address.

Notes and tips

- Use port 587 for STARTTLS or 465 for SMTPS. The sender supports both: it will use direct TLS for port 465 and STARTTLS for other ports when available.
- For Gmail SMTP, you may need an app password or OAuth; plain username/password may be blocked unless you enable app passwords.
- For self-signed SMTP servers you may need to adapt the TLS config (the current code uses strict ServerName verification). If you need a quick test with self-signed certs, I can add a config flag to allow insecure TLS.

Security

- In production, ensure TLS verification is not skipped and use a verified `SMTP_FROM` address.
- Consider using a transactional email provider for improved deliverability (SendGrid, Mailgun, SES).
