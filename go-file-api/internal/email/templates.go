package email

import (
	"bytes"
	"html/template"
)

var vaultInviteTmpl = template.Must(template.New("vault_invite").Parse(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; color: #222; max-width: 480px; margin: 40px auto;">
  <h2>You've been invited to <strong>{{.VaultName}}</strong></h2>
  <p>Someone has shared a vault with you on Go File App. Click the button below to accept the invitation and set up your account.</p>
  <a href="{{.InviteURL}}"
     style="display:inline-block;padding:12px 24px;background:#18181b;color:#fff;text-decoration:none;border-radius:6px;font-weight:600;">
    Accept Invitation
  </a>
  <p style="margin-top:32px;font-size:13px;color:#888;">
    This link expires in 7 days. If you weren't expecting this invitation, you can ignore this email.
  </p>
  <p style="font-size:13px;color:#888;">Or copy this URL into your browser:<br>{{.InviteURL}}</p>
</body>
</html>
`))

var vaultAccessGrantedTmpl = template.Must(template.New("vault_access_granted").Parse(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; color: #222; max-width: 480px; margin: 40px auto;">
  <h2>You now have access to <strong>{{.VaultName}}</strong></h2>
  <p>You've been granted access to a vault on Go File App. Log in to view its contents.</p>
</body>
</html>
`))

var welcomeTmpl = template.Must(template.New("welcome").Parse(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; color: #222; max-width: 480px; margin: 40px auto;">
  <h2>Welcome to Go File App!</h2>
  <p>Thanks for registering, <strong>{{.Email}}</strong>. Your account is ready and a default vault has been created for you.</p>
  <p>You can start uploading files and sharing vaults straight away.</p>
  <p style="margin-top:32px;font-size:13px;color:#888;">If you didn't create this account, you can ignore this email.</p>
</body>
</html>
`))

var resetPasswordTmpl = template.Must(template.New("reset_password").Parse(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; color: #222; max-width: 480px; margin: 40px auto;">
  <h2>Reset your password</h2>
  <p>We received a request to reset the password for <strong>{{.Email}}</strong>. Click the button below to choose a new one.</p>
  <a href="{{.ResetURL}}"
     style="display:inline-block;padding:12px 24px;background:#18181b;color:#fff;text-decoration:none;border-radius:6px;font-weight:600;">
    Reset Password
  </a>
  <p style="margin-top:32px;font-size:13px;color:#888;">
    This link expires in 15 minutes. If you didn't request a password reset, you can safely ignore this email.
  </p>
  <p style="font-size:13px;color:#888;">Or copy this URL into your browser:<br>{{.ResetURL}}</p>
</body>
</html>
`))

func renderVaultInvite(vaultName, inviteURL string) string {
	var buf bytes.Buffer
	vaultInviteTmpl.Execute(&buf, struct {
		VaultName string
		InviteURL string
	}{vaultName, inviteURL})
	return buf.String()
}

func renderVaultAccessGranted(vaultName string) string {
	var buf bytes.Buffer
	vaultAccessGrantedTmpl.Execute(&buf, struct{ VaultName string }{vaultName})
	return buf.String()
}

func renderWelcome(email string) string {
	var buf bytes.Buffer
	welcomeTmpl.Execute(&buf, struct{ Email string }{email})
	return buf.String()
}

func renderResetPassword(email string, resetURL string) string {
	var buf bytes.Buffer
	resetPasswordTmpl.Execute(&buf, struct {
		Email    string
		ResetURL string
	}{email, resetURL})
	return buf.String()
}
