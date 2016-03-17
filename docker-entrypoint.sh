#!/bin/bash

# 1. Run database migrations
/go/bin/pinglist-api migrate

# 2. Load fixtures
/go/bin/pinglist-api loaddata \
  oauth/fixtures/scopes.yml \
  subscriptions/fixtures/plans.yml \
  alarms/fixtures/alarm_states.yml \
  alarms/fixtures/incident_types.yml

# Finally, run the server
/go/bin/pinglist-api runserver
