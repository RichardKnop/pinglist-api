#!/bin/bash

# 1. Run database migrations
/go/bin/pinglist-api migrate

# 2. Load fixtures
/go/bin/pinglist-api loaddata \
  oauth/fixtures/scopes.yml \
  accounts/fixtures/roles.yml \
  subscriptions/fixtures/plans.yml \
  alarms/fixtures/regions.yml \
  alarms/fixtures/alarm_states.yml \
  alarms/fixtures/incident_types.yml

# Finally, run the web server and scheduler
/go/bin/pinglist-api runall
