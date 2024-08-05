FROM ghcr.io/lesomnus/dev-golang:1.22

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
	--mount=type=cache,target=/var/lib/apt,sharing=locked \
	apt update \
	&& apt-get install --no-install-recommends --yes \
		clang-format

RUN TEMP=$(mktemp -d) \
	&& cd "${TEMP}" \
	&& curl -o ./protoc.zip -sSL "https://github.com/protocolbuffers/protobuf/releases/download/v24.0/protoc-24.0-linux-x86_64.zip" \
	&& unzip ./protoc.zip -d /usr/local \
	&& chmod 755 /usr/local/bin/protoc \
	&& rm -rf "${TEMP}"

RUN TEMP=$(mktemp -d) \
	&& cd "${TEMP}" \
	&& curl -o ./a.tar.gz -sSL "https://github.com/ktr0731/evans/releases/download/v0.10.11/evans_linux_$(dpkg --print-architecture).tar.gz" \
	&& tar -xf ./a.tar.gz \
	&& mv ./evans /usr/local/bin/. \
	&& rm -rf "${TEMP}" \
	&& evans --version

RUN TEMP=$(mktemp -d) \
	&& cd "${TEMP}" \
	&& curl -o ./a -L "https://dl.k8s.io/release/$(curl -sSL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
	&& chmod +x ./a \
	&& mv ./a /usr/local/bin/kubectl \
	&& rm -rf "${TEMP}" \
	&& kubectl version --client

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
	--mount=type=cache,target=/var/lib/apt,sharing=locked \
	curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc \
	&& chmod a+r /etc/apt/keyrings/docker.asc \
	&& echo \
		"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
		$(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
		| tee /etc/apt/sources.list.d/docker.list > /dev/null \
	&& apt update \
	&& apt-get install --no-install-recommends --yes \
		docker-ce-cli \
	&& docker --version

RUN mkdir -p /usr/local/include/google/api/ \
	&& cd /usr/local/include/google/api/ \
	&& curl -sSLO "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/{annotations,http}.proto"



USER hypnos

RUN docker context create dev --docker "host=tcp://docker:2376,skip-tls-verify=true" \
	&& docker context use dev

RUN  echo \
		google.golang.org/protobuf/cmd/protoc-gen-go@v1.31 \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3 \
		entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc@master \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.21.0 \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.21.0 \
		github.com/yoheimuta/protolint/cmd/protolint@latest \
	| xargs -n1 go install

# Comment below lines if you don't want to use Zsh
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true \
	&& git clone https://github.com/zsh-users/zsh-autosuggestions         ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions     \
	&& git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting \
	&& git clone https://github.com/zsh-users/zsh-completions             ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions         \
	&& git clone https://github.com/supercrabtree/k                       ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/k
COPY --chown=1000:1000 .zshrc /home/hypnos
ENV SHELL /bin/zsh
