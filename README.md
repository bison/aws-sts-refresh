# aws-sts-refresh

Refresh saved AWS STS credentials using the existing (hopefully not
expired) credentials.

```
aws-sts-refresh -p 'my-profile' -r 'arn:aws:iam::123456789:role/some-role' -n 'my-name'
```
