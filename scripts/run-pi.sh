#!/bin/sh
/home/j/nasblaze/nasblaze \
    -log_dir=/home/j/logs \
    --dry-run=false \
    --hard-delete=false \
    --exclude="**/node_modules/**" \
    --exclude="**/.next/**" \
    --exclude="**/__pycache__/**" \
    --exclude="**/.git/**" \
    --exclude="**/.svn/**" \
    --exclude="**/.metadata/**" \
    --exclude="Programming/travive/**" \
    --exclude=".Spotlight-V100/**" \
    --exclude=".com.apple.timemachine.donotpresent" \
    --exclude="._.com.apple.timemachine.donotpresent" \
    --exclude="**/*.class" \
    --exclude=".DS_Store" \
    --exclude="._.DS_Store" \
    --exclude="_DS_Store" \
    --exclude="do_not_backup/**"
