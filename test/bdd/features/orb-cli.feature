
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@all
@orb_cli
Feature: Using Orb CLI
  Background: Setup

    Given domain "orb.domain1.com" is mapped to "localhost:48326"
    And domain "orb.domain2.com" is mapped to "localhost:48426"

    Given the authorization bearer token for "GET" requests to path "/sidetree/v1/identifiers" is set to "READ_TOKEN"
    And the authorization bearer token for "POST" requests to path "/log" is set to "ADMIN_TOKEN"

    # set up logs for domains
    When an HTTP POST is sent to "https://orb.domain1.com/log" with content "http://orb.vct:8077/maple2020" of type "text/plain"

    And orb-cli is executed with args 'acceptlist add --url https://localhost:48326/services/orb/acceptlist --actor https://orb.domain2.com/services/orb --type follow --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'
    And orb-cli is executed with args 'acceptlist add --url https://localhost:48426/services/orb/acceptlist --actor https://orb.domain1.com/services/orb --type invite-witness --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'

  @orb_cli_did
  Scenario: Create and update did doc using CLI and verify proofs in VCT log
    # domain2 server follows domain1 server
    When user create "follower" activity with outbox-url "https://localhost:48426/services/orb/outbox" actor "https://orb.domain2.com/services/orb" to "https://orb.domain1.com/services/orb" action "Follow"
    # domain1 invites domain2 to be a witness
    When user create "witness" activity with outbox-url "https://localhost:48326/services/orb/outbox" actor "https://orb.domain1.com/services/orb" to "https://orb.domain2.com/services/orb" action "InviteWitness"
    Then we wait 3 seconds
    When Create keys in kms
    When Orb DID is created through cli
    Then check cli created valid DID
    And we wait 3 seconds

    When Orb DID is resolved through cli
    Then the JSON path "didDocumentMetadata.versionId" of the response is saved to variable "anchorHash"
    And we wait 3 seconds

    When orb-cli is executed with args 'vct verify --cas-url https://localhost:48326/cas --anchor ${anchorHash} --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN --vct-auth-token=vctread'
    Then the JSON path '#(domain=="http://orb.vct:8077/maple2020").found' of the boolean response equals "true"

    When Orb DID is updated through cli
    Then check cli updated DID
    When Orb DID is recovered through cli
    Then check cli recovered DID
    When Orb DID is deactivated through cli
    Then check cli deactivated DID

  @orb_cli_activity
  Scenario: test follow and witness
    # domain1 server follows domain2 server
    When user create "follower" activity with outbox-url "https://localhost:48326/services/orb/outbox" actor "https://orb.domain1.com/services/orb" to "https://orb.domain2.com/services/orb" action "Follow"
    Then we wait 3 seconds
    When user create "follower" activity with outbox-url "https://localhost:48326/services/orb/outbox" actor "https://orb.domain1.com/services/orb" to "https://orb.domain2.com/services/orb" action "Undo"

      # domain2 invites domain1 to be a witness
    When user create "witness" activity with outbox-url "https://localhost:48326/services/orb/outbox" actor "https://orb.domain1.com/services/orb" to "https://orb.domain2.com/services/orb" action "InviteWitness"
    Then we wait 3 seconds
    When user create "witness" activity with outbox-url "https://localhost:48326/services/orb/outbox" actor "https://orb.domain1.com/services/orb" to "https://orb.domain2.com/services/orb" action "Undo"

  @orb_cli_acceptlist
  Scenario: test accept list management using cli
    # Add actors to the 'follow' accept list.
    When orb-cli is executed with args 'acceptlist add --url https://localhost:48326/services/orb/acceptlist --actor https://orb.domainx.com/services/orb --actor https://orb.domainy.com/services/orb --type follow --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'
    # Add actors to the 'invite-witness' accept list.
    Then orb-cli is executed with args 'acceptlist add --url https://localhost:48326/services/orb/acceptlist --actor https://orb.domainz.com/services/orb --type invite-witness --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'

    When orb-cli is executed with args 'acceptlist get --url https://localhost:48326/services/orb/acceptlist --type follow --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the JSON path "url" of the response contains "https://orb.domainx.com/services/orb"
    Then the JSON path "url" of the response contains "https://orb.domainy.com/services/orb"

    And orb-cli is executed with args 'acceptlist get --url https://localhost:48326/services/orb/acceptlist --type invite-witness --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the JSON path "url" of the response contains "https://orb.domainz.com/services/orb"

    And orb-cli is executed with args 'acceptlist get --url https://localhost:48326/services/orb/acceptlist --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the JSON path '#(type="follow").url' of the response contains "https://orb.domainx.com/services/orb"
    And the JSON path '#(type="follow").url' of the response contains "https://orb.domainx.com/services/orb"
    And the JSON path '#(type="invite-witness").url' of the response contains "https://orb.domainz.com/services/orb"

    # Remove actors from the 'follow' accept list.
    When orb-cli is executed with args 'acceptlist remove --url https://localhost:48326/services/orb/acceptlist --actor https://orb.domainx.com/services/orb --actor https://orb.domainy.com/services/orb --type follow --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'
    # Remove actors from the 'invite-witness' accept list.
    Then orb-cli is executed with args 'acceptlist remove --url https://localhost:48326/services/orb/acceptlist --actor https://orb.domainz.com/services/orb --type invite-witness --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'

    And orb-cli is executed with args 'acceptlist get --url https://localhost:48326/services/orb/acceptlist --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the JSON path '#(type="follow").url' of the response does not contain "https://orb.domainx.com/services/orb"
    And the JSON path '#(type="follow").url' of the response does not contain "https://orb.domainx.com/services/orb"
    And the JSON path '#(type="invite-witness").url' of the response does not contain "https://orb.domainz.com/services/orb"

  @orb_cli_policy
  Scenario: test witness policy using cli
    When orb-cli is executed with args 'policy update --url https://localhost:48326/policy --policy "MinPercent(100,batch) AND OutOf(1,system)" --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'
    And orb-cli is executed with args 'policy get --url https://localhost:48326/policy --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the response equals "MinPercent(100,batch) AND OutOf(1,system)"

  @orb_cli_log
  Scenario: test domain log management using cli
    When orb-cli is executed with args 'log update --url https://localhost:48326/log --log "http://orb.vct:8077/maple2022" --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'
    And orb-cli is executed with args 'policy get --url https://localhost:48326/log --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the response equals "http://orb.vct:8077/maple2022"

  @orb_cli_logmonitor
  Scenario: test log monitor management using cli
    # Add log for monitoring.
    When orb-cli is executed with args 'logmonitor activate --url https://localhost:48326/log-monitor --log http://orb.vct:8077/maple2022 --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'

    Then we wait 1 seconds

    When orb-cli is executed with args 'logmonitor get --url https://localhost:48326/log-monitor --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the JSON path "active.#.log_url" of the response contains "http://orb.vct:8077/maple2022"

     # Deactivate log - remove from log monitoring list.
    When orb-cli is executed with args 'logmonitor deactivate --url https://localhost:48326/log-monitor --log http://orb.vct:8077/maple2022 --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token ADMIN_TOKEN'

    Then we wait 1 seconds

    When orb-cli is executed with args 'logmonitor get --url https://localhost:48326/log-monitor --status inactive --tls-cacerts fixtures/keys/tls/ec-cacert.pem --auth-token READ_TOKEN'
    Then the JSON path "inactive.#.log_url" of the response contains "http://orb.vct:8077/maple2022"

