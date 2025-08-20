#!/bin/bash

# ───────────────────────────────────────────────────────────
# oci_auth_refresher.sh  •  v0.2.2
#
# Keeps an OCI CLI session alive by refreshing it shortly
# before it expires. Intended to be launched nohup
# ───────────────────────────────────────────────────────────

# Check if OCI CLI is installed
if ! command -v oci >/dev/null 2>&1; then
  echo "[ERROR] OCI CLI not found. Please install it and ensure it's in your PATH."
  exit 1
fi

# Resolve profile
if [[ -n "$1" ]]; then
  OCI_CLI_PROFILE=$1
elif [[ -n "${OCI_CLI_PROFILE}" ]]; then
  :
else
  OCI_CLI_PROFILE="DEFAULT"
fi

# Session directory
SESSION_DIR="${HOME}/.oci/sessions/${OCI_CLI_PROFILE}"
SESSION_STATUS_FILE="${SESSION_DIR}/session_status"
REFRESHER_PID_FILE="${SESSION_DIR}/refresher.pid"
LOG_FILE="${SESSION_DIR}/refresher.log"

mkdir -p "$SESSION_DIR"

# Relaunch with nohup
if [[ -z "$NOHUP" && -t 1 ]]; then
  export NOHUP=1
  script_path="$(cd "$(dirname "$0")" && pwd)/$(basename "$0")"
  nohup "$script_path" "$OCI_CLI_PROFILE" > "$LOG_FILE" 2>&1 < /dev/null &
  exit 0
fi

PREEMPT_REFRESH_TIME=60 # Attempt to refresh 60 sec before session expiration

# Log session expired
function log_session_expired() {
  local reason="$1"
  oci_session_status="expired"
  echo "$oci_session_status" > "$SESSION_STATUS_FILE"
  echo "[INFO] Session expired. Reason: ${reason}" >> "$LOG_FILE"
}

function cleanup_pid_file() {
  rm -f "$REFRESHER_PID_FILE"
}
trap cleanup_pid_file EXIT

# Convert timestamp to epoch
function to_epoch() {
  local ts="$1"
  [[ -z "$ts" ]] && return 1

  if date --version >/dev/null 2>&1; then
    date -d "${ts}" +%s 2>/dev/null || return 1
  else
    for fmt in "%Y-%m-%d %H:%M:%S" "%Y-%m-%d %T" "%Y-%m-%dT%H:%M:%S" "%Y-%m-%d"; do
      if date -j -f "$fmt" "${ts}" +%s 2>/dev/null; then
        return 0
      fi
    done
    return 1
  fi
}

# Check remaining session duration
function get_remaining_session_duration() {
  if oci session validate --profile "$OCI_CLI_PROFILE" --local >> "$LOG_FILE" 2>&1; then
    oci_session_status="valid"
    echo "$oci_session_status" > "$SESSION_STATUS_FILE"

    local validate_output
    validate_output=$(oci session validate --profile "$OCI_CLI_PROFILE" --local 2>&1)

    local exp_ts
    exp_ts=$(echo "$validate_output" | sed -E 's/.*until ([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}).*/\1/')
    [[ "$exp_ts" == "$validate_output" ]] && exp_ts=""

    if [[ -z "$exp_ts" ]]; then
      local date_part time_part
      date_part=$(echo "$validate_output" | grep -o "[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}")
      time_part=$(echo "$validate_output" | grep -o "[0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\}")
      [[ -n "$date_part" && -n "$time_part" ]] && exp_ts="$date_part $time_part"
    fi

    if [[ -z "$exp_ts" || ! "$exp_ts" =~ [0-9]{4}-[0-9]{2}-[0-9]{2} ]]; then
      log_session_expired "invalid timestamp"
      remaining_time=0
      return
    fi

    local exp_epoch
    if ! exp_epoch=$(to_epoch "${exp_ts}"); then
      log_session_expired "invalid epoch conversion"
      remaining_time=0
      return
    fi

    local now_epoch
    now_epoch=$(date +%s)
    remaining_time=$((exp_epoch - now_epoch))
  else
    log_session_expired ""
    remaining_time=0
  fi
}

function refresh_session() {
  if oci session refresh --profile "$OCI_CLI_PROFILE" 2>&1 >> "$LOG_FILE"; then
    echo "[INFO] Session refreshed successfully." >> "$LOG_FILE"
    return 0
  else
    log_session_expired "refresh failed"
    return 1
  fi
}

# Init
oci_session_status="unknown"
remaining_time=0

get_remaining_session_duration

while [[ "$oci_session_status" == "valid" ]]; do
  if (( remaining_time > PREEMPT_REFRESH_TIME )); then
    sleep_for=$((remaining_time - PREEMPT_REFRESH_TIME))
    echo "[INFO] Sleeping for $sleep_for seconds before next refresh..." >> "$LOG_FILE"
    sleep "$sleep_for"

    if ! refresh_session; then
      exit 1
    fi

    get_remaining_session_duration
  else
    log_session_expired "main loop"
  fi
done

exit 0