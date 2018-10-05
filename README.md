# BucketServer
Serve GCS Bucket content anonymously

Serves a private bucket on Google Cloud Storage for anonymouse users
by redirecting to a signed url of an object. A key file for a service
user needs to be created. BucketServer will bind the specified port to
all IP addresses on the local node. Please make sure you do not serve
your private bucket to an unintended network. For example a VM on GCP
with an external address might lead to serving the private bucket to the
internet.

Disclaimer: This is not an officially supported Google product

Example:

```powershell
# Install BucketServer
go get -u -v github.com/google/BucketServer
# Create a service account
$AccountName = 'bucketserver-sa'
gcloud iam service-accounts create $AccountName --display-name "Service account for BucketServer"
# Get the Email of the service account
$Email = (gcloud iam service-accounts list --filter "EMAIL:($AccountName@*)")[1].Split() | Select-Object -Last 1
# Create a key JSON file for this service account
gcloud iam service-accounts keys create key.json --iam-account $Email
# Get a unique bucket name
$BucketName = [System.Guid]::NewGuid()
# Create the bucket
gsutil mb -c regional -l europe-west4 gs://$BucketName
# Grant read access to the service account
gsutil iam ch -e '' -d serviceAccount:$($Email):objectViewer gs://$BucketName
# Copy File to bucket
Get-Date > t.txt
gsutil cp t.txt gs://$BucketName
# Run BucketServer
$key = Get-Childitem key.json
$job = Start-Job -ScriptBlock {  &"$ENV:GOPATH\bin\BucketServer" $Args[0] $Args[1] $Args[2] } -ArgumentList @($BucketName, $key.Fullname, "8080")
# Should be running
$job | Get-Job
# Request object
Invoke-WebRequest -Uri http://localhost:8080/t.txt
# Output
# StatusCode        : 200
# StatusDescription : OK
# Content           : ÿþ
#                      S a t u r d a y ,   8   S e p t e m b e r ,   2 0 1 8   0 4 : 3 6 : 3 4
# Stop Server
$job | Stop-Job
$job | Remove-Job
```
