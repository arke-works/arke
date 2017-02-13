FROM       scratch
MAINTAINER Aleksandr Tihomirov <hunter@zeta.pm>
ADD        arke-forum arke-forum
ENV        PORT 80
EXPOSE     80
ENTRYPOINT ["/arke-forum"]