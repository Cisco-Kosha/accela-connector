# Kosh Accela Connector

Accela is a cloud-based platform with which you can build and deploy applications and services to streamline and automate government processes, such as permitting, licensing, and code enforcement.

The Kosha Accela connector enables you to perform REST API operations from the Accela API in your Kosha workflow or custom application. Using the Kosha Accela connector, you can directly access the Accela platform to:

* Manage an agency's records
* Update a record's activities
* Manage citizen access users
* Manage assessments

## Useful Actions

You can use the Kosha Accela connector to manage records, citizen access user management, assessment, and assets.

Refer to the Kosha Accela connector [API specification](openapi.json) for details.

### Records

Use the records API to:

* Create, list, update, and delete records
* Create, list, update, and delete record comments
* Create, list, update, and delete conditional approvals
* Create and estimate record fees
  
### Citizen Access User Management

Use the citizen access user management API to:

* Create citizen users
* Add and delete citizen contacts
* Update citizen passwords and account statements

### Assessments

Use the assessments API to:

* Get, delete, and update assessments
* Create and delete assessment documents

### Assets

Use the assets API to:

* Create, delete, and update assets
* Create, attach, and delete asset documents

## Authentication

To authenticate when provisioning the Kosha Accela connector, you need your:

* Application client ID and client secret
* Accela account username and password
* Accela environment name, such as PROD or TEST
* Permission scopes needed to request from the Accela API server

## Kosha Connector Open Source Development

All connectors Kosha shares on the marketplace are open source. We believe in fostering collaboration and open development. Everyone is welcome to contribute their ideas, improvements, and feedback for any Kosha connector. We encourage community engagement and appreciate any contributions that align with our goals of an open and collaborative API management platform.

Refer to the contribution guidelines for details.

## Contributing

Pull requests and bug reports are welcome.

For larger changes, please create an issue in GitHub first to discuss your proposed changes and their possible implications.
