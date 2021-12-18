#!/bin/sh

EXTRA_ARGS=$EXTRA_ARGS
if [ $LISTENPORT ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ListenPort='$LISTENPORT
fi

if [ $ENDPOINTNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointName='$ENDPOINTNAME
fi

if [ $DBSERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -DBServer='$DBSERVER
fi

if [ $DBNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -DBName='$DBNAME
fi

if [ $DBUSERNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -DBUsername='$DBUSERNAME
fi

if [ $DBPASSWORD ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -DBPassword='$DBPASSWORD
fi

if [ $BATISSERVICE ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -BatisService='$BATISSERVICE
fi

if [ $CASSERVICE ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -CasService='$CASSERVICE
fi

if [ $FILESERVICE ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -FileService='$FILESERVICE
fi

if [ $SUPERNAMESPACE ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -SuperNamespace='$SUPERNAMESPACE
fi

if [ $DEV ]; then
    echo $EXTRA_ARGS
fi

/var/app/magicDefault $EXTRA_ARGS "$@"
