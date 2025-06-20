ARG ALPINE_CONTAINER_IMAGE=alpine:3.20.3

#=========================
# Build the app container.
# Patch the alpine image.
#=========================
FROM ${ALPINE_CONTAINER_IMAGE} AS appcontainer

# Add unprivileged user. "app" is the user that will run the app.
RUN addgroup -S app
RUN adduser -S -G app app

# Remove unnecessary accounts
RUN sed -i -r "/^(root|nobody|app)/!d" /etc/group \
    && sed -i -r "/^(root|nobody|app)/!d" /etc/passwd

# Remove init scripts since we do not use them.
RUN rm -fr /etc/init.d /lib/rc /etc/conf.d /etc/inittab /etc/runlevels /etc/rc.conf /etc/logrotate.d

# Remove root home dir
RUN rm -fr /root

# Remove any symlinks that we broke during previous steps
RUN find /bin /etc /lib /sbin /usr -xdev -type l -exec test ! -e {} \; -delete

#===========================================================
# Base container for gathering files and setting permissions
#===========================================================
FROM ${ALPINE_CONTAINER_IMAGE} AS gather-files-base

# Add unprivileged user. "app" is the user that will run the app.
RUN addgroup -S app
RUN adduser -S -G app app

#=======================================================
# EXPLODING KITTENS - Gather, set permissions, and build the kittens image
#=======================================================
FROM gather-files-base AS gather-kittens
COPY cmd/exploding-kittens/exploding-kittens /cmd/exploding-kittens/exploding-kittens
RUN chown -R app:app /cmd/exploding-kittens/exploding-kittens

FROM appcontainer AS kittens
COPY --from=gather-kittens /cmd/exploding-kittens/exploding-kittens /app
USER app
CMD [ "./app" ]
