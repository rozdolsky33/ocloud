#!/bin/zsh
# shellcheck shell=bash disable=SC1071

# ───────────────────────────────────────────────────────────
# oci_auth_refresher.sh  •  v0.0.21
#
# Keeps an OCI CLI session alive by refreshing it shortly
# before it expires. Intended to be launched (nohup) from the
# wrapper script oshell.sh.
# ───────────────────────────────────────────────────────────

# Check if profile argument is provided, use DEFAULT if not
if [[ -z "$1" ]]; then
  echo "No profile name provided, using DEFAULT"
  OCI_CLI_PROFILE="DEFAULT"
else
  OCI_CLI_PROFILE=$1
fi

# Check if script is being run directly (not through nohup)
# If so, relaunch itself using nohup and exit
if [[ -z "$NOHUP" && -t 1 ]]; then
  echo "Launching OCI auth refresher in background for profile ${OCI_CLI_PROFILE}"
  export NOHUP=1
  # Use full path to script to ensure it's detectable by pgrep
  script_path=$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)/$(basename "$0")
  nohup "$script_path" "$OCI_CLI_PROFILE" > /dev/null 2>&1 < /dev/null &
  pid=$!
  echo "Process started with PID $pid"
  exit 0
fi
  # Configuration
  PREEMPT_REFRESH_TIME=60  # Attempt to refresh 60 sec before session expiration
  SESSION_STATUS_FILE="${HOME}/.oci/sessions/${OCI_CLI_PROFILE}/session_status"


# Function to get the remaining duration of the current session
function get_remaining_session_duration() {

  if oci session validate --profile "$OCI_CLI_PROFILE" --local 2>&1; then
    oci_session_status="valid"
    echo "$oci_session_status" > "$SESSION_STATUS_FILE"

    # Get expiration timestamp
    local exp_ts
    local validate_output

    # Capture both stdout and stderr
    validate_output=$(oci session validate --profile "$OCI_CLI_PROFILE" --local 2>&1)

    # Use a simple approach to extract the date and time
    exp_ts=$(echo "$validate_output" | sed -E 's/.*until ([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}).*/\1/')

    # If the output is unchanged, it means the pattern didn't match
    if [[ "$exp_ts" == "$validate_output" ]]; then
      exp_ts=""
    fi

    # If still empty, try to extract just the date and time parts
    if [[ -z "$exp_ts" ]]; then
      local date_part
      date_part=$(echo "$validate_output" | grep -o "[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}")
      local time_part
      time_part=$(echo "$validate_output" | grep -o "[0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\}")

      if [[ -n "$date_part" && -n "$time_part" ]]; then
        exp_ts="$date_part $time_part"
    fi

    # Verify that we have a valid-looking timestamp before proceeding
    if [[ -z "$exp_ts" || ! "$exp_ts" =~ [0-9]{4}-[0-9]{2}-[0-9]{2} ]]; then
      oci_session_status="expired"
      echo "$oci_session_status" > "$SESSION_STATUS_FILE"
      remaining_time=0
      return
    fi

    # Calculate remaining time
    local exp_epoch
    if ! exp_epoch=$(to_epoch "${exp_ts}"); then
      oci_session_status="expired"
      echo "$oci_session_status" > "$SESSION_STATUS_FILE"
      remaining_time=0
      return
    fi
    local now_epoch
    now_epoch=$(date +%s)
    remaining_time=$((exp_epoch - now_epoch))
  else
    oci_session_status="expired"
    echo "$oci_session_status" > "$SESSION_STATUS_FILE"
    remaining_time=0
  fi
}

# Function to refresh the session
function refresh_session() {

  if oci session refresh --profile "$OCI_CLI_PROFILE" 2>&1; then
    return 0
  else
    oci_session_status="expired"
    echo "$oci_session_status" > "$SESSION_STATUS_FILE"
    return 1
  fi
}

# Initialize variables
oci_session_status="unknown"
remaining_time=0

# Check if session directory exists
if [[ ! -d "${HOME}/.oci/sessions/${OCI_CLI_PROFILE}" ]]; then
  echo "Missing session directory; user probably hasn't authenticated"
  echo "Exiting."
  exit 1
fi

# Main loop
get_remaining_session_duration

while [[ "$oci_session_status" == "valid" ]]; do
  if (( remaining_time > PREEMPT_REFRESH_TIME )); then
    sleep_for=$((remaining_time - PREEMPT_REFRESH_TIME))
    sleep "$sleep_for"

    if ! refresh_session; then
      echo "Exiting due to refresh failure"
      exit 1
    fi

    get_remaining_session_duration
  else
    oci_session_status="expired"
    echo "$oci_session_status" > "$SESSION_STATUS_FILE"
  fi
done

exit 0
