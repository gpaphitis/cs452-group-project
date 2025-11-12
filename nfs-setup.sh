#!/bin/bash
# setup-nfs.sh
# Usage:
#   Server: sudo ./setup-nfs.sh server [SUBNET_CIDR]
#           e.g. sudo ./setup-nfs.sh server 10.10.1.0/24
#           (omit SUBNET_CIDR to auto-detect a 10.x.x.x/24)
#   Client: sudo ./setup-nfs.sh client <SERVER_HOST_OR_IP>
#
# Notes:
# - Server exports: /srv/mapReduceData
# - Clients mount : /mnt/mapReduceData
# - Idempotent-ish: avoids duplicate /etc/exports and /etc/fstab lines

set -euo pipefail

ROLE="${1:-}"
SUBNET_CIDR="${2:-}"
SERVER_HOST="${2:-}" # reused for client role

SHARE_DIR="/srv/mapReduceData"
MOUNT_POINT="/mnt/mapReduceData"
EXPORT_OPTS="rw,sync,no_subtree_check,no_root_squash"
FSTAB_OPTS="nfs defaults,_netdev 0 0"

need_root() {
  if [[ $EUID -ne 0 ]]; then
    echo "Please run as root (use sudo)." >&2
    exit 1
  fi
}

dedupe_line_in_file() {
  local line="$1" file="$2"
  # remove any exact duplicates of the line, then append once
  if [[ -f "$file" ]]; then
    grep -Fxv "$line" "$file" > "${file}.tmp" || true
    mv "${file}.tmp" "$file"
  fi
  echo "$line" >> "$file"
}

detect_subnet_cidr() {
  # Try to auto-detect a 10.x.x.x/24 network (common on CloudLab experiment LANs)
  local subnet
  subnet=$(ip -4 addr show | grep -Eo '10\.[0-9]+\.[0-9]+\.0/24' | head -n1 || true)
  if [[ -z "$subnet" ]]; then
    # Fallback: derive /24 from a 10.x.x.x address if available
    local ip
    ip=$(ip -4 addr show | grep -Eo '10\.[0-9]+\.[0-9]+\.[0-9]+' | head -n1 || true)
    if [[ -n "$ip" ]]; then
      subnet="${ip%.*}.0/24"
    fi
  fi
  echo "$subnet"
}

setup_server() {
  need_root

  apt-get update -y
  apt-get install -y nfs-kernel-server

  mkdir -p "$SHARE_DIR"
  chmod 777 "$SHARE_DIR"

  local cidr="${SUBNET_CIDR}"
  if [[ -z "$cidr" ]]; then
    cidr=$(detect_subnet_cidr || true)
  fi

  if [[ -z "$cidr" ]]; then
    echo "Could not auto-detect a 10.x.x.x/24 subnet."
    echo "    Re-run with an explicit CIDR, e.g.:"
    echo "    sudo $0 server 10.10.1.0/24"
    exit 1
  fi

  local export_line="${SHARE_DIR} ${cidr}(${EXPORT_OPTS})"
  touch /etc/exports
  # Remove any existing lines that export SHARE_DIR (to any net), then add our fresh one
  grep -vE "^${SHARE_DIR}\s" /etc/exports > /etc/exports.tmp || true
  mv /etc/exports.tmp /etc/exports
  dedupe_line_in_file "$export_line" /etc/exports

  exportfs -ra
  systemctl enable --now nfs-server

  echo "✅ NFS server ready."
  echo "   Exported: ${SHARE_DIR}"
  echo "   Allowed:  ${cidr}"
  echo
  showmount -e localhost || true
}

setup_client() {
  need_root

  if [[ -z "$SERVER_HOST" ]]; then
    echo "Usage: sudo $0 client <SERVER_HOST_OR_IP>" >&2
    exit 1
  fi

  apt-get update -y
  apt-get install -y nfs-common

  mkdir -p "$MOUNT_POINT"

  # Try NFSv4 first; fall back to default if needed
  set +e
  mount -t nfs4 "${SERVER_HOST}:${SHARE_DIR}" "$MOUNT_POINT" 2>/dev/null
  rc=$?
  set -e
  if [[ $rc -ne 0 ]]; then
    mount "${SERVER_HOST}:${SHARE_DIR}" "$MOUNT_POINT"
  fi

  # Add to /etc/fstab if not present (dedupe exact line)
  local fstab_line="${SERVER_HOST}:${SHARE_DIR} ${MOUNT_POINT} ${FSTAB_OPTS}"
  touch /etc/fstab
  grep -vF "${SERVER_HOST}:${SHARE_DIR} ${MOUNT_POINT} " /etc/fstab > /etc/fstab.tmp || true
  mv /etc/fstab.tmp /etc/fstab
  dedupe_line_in_file "$fstab_line" /etc/fstab

  echo "Client mounted ${SERVER_HOST}:${SHARE_DIR} at ${MOUNT_POINT}"
}

case "$ROLE" in
  server)
    setup_server
    ;;
  client)
    setup_client
    ;;
  *)
    cat >&2 <<EOF
Usage:
  Server: sudo $0 server [SUBNET_CIDR]
          e.g. sudo $0 server 10.10.1.0/24
  Client: sudo $0 client <SERVER_HOST_OR_IP>
          e.g. sudo $0 client node0
EOF
    exit 1
    ;;
esac