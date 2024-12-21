FROM hashicorp/consul:latest

EXPOSE 8500 8600/udp

CMD ["agent", "-server", "-ui", "-node=server-1", "-bootstrap-expect=1", "-client=0.0.0.0"]
