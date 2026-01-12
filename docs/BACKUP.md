# Database Backup Guide

## Overview

The worklog application includes automated PostgreSQL database backups to AWS S3. Backups run daily at 02:00 UTC and are retained for 90 days.

## Setup Instructions

### 1. Create S3 Bucket

Create an S3 bucket for storing backups:

```bash
aws s3 mb s3://your-backup-bucket-name
```

### 2. Configure S3 Lifecycle Policy

Set up automatic deletion of backups older than 90 days:

```bash
cat > lifecycle-policy.json << 'EOF'
{
  "Rules": [
    {
      "Id": "DeleteOldBackups",
      "Status": "Enabled",
      "Expiration": {
        "Days": 90
      }
    }
  ]
}
EOF

aws s3api put-bucket-lifecycle-configuration \
  --bucket worklog-backups \
  --lifecycle-configuration file://lifecycle-policy.json
```

### 3. Create IAM User for Backups

Create an IAM user with permissions to upload to the S3 bucket:

```bash
cat > backup-policy.json << 'EOF'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::your-backup-bucket-name",
        "arn:aws:s3:::your-backup-bucket-name/*"
      ]
    }
  ]
}
EOF

# Create IAM user
aws iam create-user --user-name worklog-backup

# Attach policy
aws iam put-user-policy \
  --user-name worklog-backup \
  --policy-name WorklogBackupPolicy \
  --policy-document file://backup-policy.json

# Create access keys
aws iam create-access-key --user-name worklog-backup
```

Save the Access Key ID and Secret Access Key from the output.

### 4. Configure GitHub Secrets

Add the following secrets to your GitHub repository (Settings > Secrets and variables > Actions):

- `AWS_ACCESS_KEY_ID` - IAM user access key ID
- `AWS_SECRET_ACCESS_KEY` - IAM user secret access key
- `AWS_REGION` - AWS region (e.g., `eu-north-1`, `us-east-1`)
- `WORKLOG_BACKUP_S3_BUCKET` - S3 bucket name (e.g., `worklog-backups`)

Note: The workflow uses the existing `WORKLOG_POSTGRES_PASSWORD` secret for database access.

## Backup Details

### Schedule

- Runs daily at 02:00 UTC
- Can be triggered manually via GitHub Actions UI

### Backup Format

- Format: PostgreSQL custom dump format (compressed)
- Compression: Level 9 (maximum)
- Naming: `worklog-backup-{github-run-id}-{timestamp}.dump`
- Example: `worklog-backup-123456789-20260112-020000.dump`

### S3 Storage Class

Backups use `STANDARD_IA` (Infrequent Access) storage class to reduce costs while maintaining quick access when needed.

### Retention

Backups are automatically deleted after 90 days via S3 lifecycle policy.

## Manual Operations

### Trigger Backup Manually

1. Go to GitHub Actions tab
2. Select "Database Backup" workflow
3. Click "Run workflow"
4. Select branch and click "Run workflow"

### List Available Backups

```bash
aws s3 ls s3://worklog-backups/ --human-readable
```

### Download a Backup

```bash
aws s3 cp s3://worklog-backups/worklog-backup-123456789-20260112-020000.dump ./backup.dump
```

## Restore Procedure

### Option 1: Restore to Kubernetes Database

```bash
# Download the backup
aws s3 cp s3://worklog-backups/worklog-backup-123456789-20260112-020000.dump ./backup.dump

# Get database password
POSTGRES_PASSWORD=$(kubectl get secret postgres-secret -n worklog -o jsonpath='{.data.POSTGRES_PASSWORD}' | base64 -d)

# Copy backup to postgres pod
kubectl cp ./backup.dump worklog/$(kubectl get pod -n worklog -l app=postgres -o jsonpath='{.items[0].metadata.name}'):/tmp/backup.dump

# Restore database (WARNING: This will overwrite existing data!)
kubectl exec -n worklog deployment/postgres -- env \
  PGPASSWORD="${POSTGRES_PASSWORD}" \
  pg_restore -U worklog -d worklog --clean --if-exists /tmp/backup.dump

# Clean up
kubectl exec -n worklog deployment/postgres -- rm /tmp/backup.dump
rm ./backup.dump
```

### Option 2: Restore to Local Database

```bash
# Download the backup
aws s3 cp s3://worklog-backups/worklog-backup-123456789-20260112-020000.dump ./backup.dump

# Restore to local database (WARNING: This will overwrite existing data!)
PGPASSWORD=yourpassword pg_restore -h localhost -U worklog -d worklog --clean --if-exists ./backup.dump

# Clean up
rm ./backup.dump
```

## Monitoring

### Check Backup Status

1. Go to GitHub Actions tab
2. View "Database Backup" workflow runs
3. Check the summary for backup details:
   - Database name
   - Timestamp
   - File name
   - File size
   - S3 location
   - Retention period

### Verify Backup in S3

```bash
# List recent backups
aws s3 ls s3://worklog-backups/ --human-readable --recursive

# Get backup metadata
aws s3api head-object \
  --bucket worklog-backups \
  --key worklog-backup-123456789-20260112-020000.dump
```

## Troubleshooting

### Backup Fails with "Access Denied"

- Verify IAM user has correct permissions
- Check that AWS credentials in GitHub secrets are correct
- Ensure S3 bucket name is correct

### Backup File is Empty

- Check that PostgreSQL pod is running: `kubectl get pods -n worklog`
- Verify database credentials are correct
- Check PostgreSQL logs: `kubectl logs -n worklog deployment/postgres`

### Cannot Restore Backup

- Verify backup file is not corrupted: `pg_restore --list backup.dump`
- Ensure PostgreSQL version matches (currently using version 16)
- Check database exists and user has permissions

## Cost Estimation

### S3 Storage Costs (eu-north-1 example)

- STANDARD_IA: $0.0125 per GB per month
- Typical backup size: ~1-10 MB (depending on data)
- 90 days of daily backups: ~90 files
- Estimated monthly cost: Less than $0.50

### Cost Optimization Tips

1. Use S3 Intelligent-Tiering for automatic cost optimization
2. Enable S3 versioning only if needed (adds cost)
3. Consider weekly backups instead of daily if acceptable
4. Reduce retention period if 90 days is not required

## Security Best Practices

1. Enable S3 bucket encryption at rest
2. Use least-privilege IAM policies
3. Rotate AWS access keys regularly
4. Enable S3 access logging for audit trail
5. Consider using AWS KMS for encryption key management
6. Restrict S3 bucket access with bucket policies

## Additional Resources

- [PostgreSQL pg_dump Documentation](https://www.postgresql.org/docs/current/app-pgdump.html)
- [PostgreSQL pg_restore Documentation](https://www.postgresql.org/docs/current/app-pgrestore.html)
- [AWS S3 Lifecycle Configuration](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lifecycle-mgmt.html)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
