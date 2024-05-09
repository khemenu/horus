FROM ghcr.io/lesomnus/dev-golang:1.21

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
	--mount=type=cache,target=/var/lib/apt,sharing=locked \
	apt update \
	&& apt-get install --no-install-recommends --yes \
		clang-format

RUN TEMP=$(mktemp -d) \
	&& cd "${TEMP}" \
	&& curl -o protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v24.0/protoc-24.0-linux-x86_64.zip \
	&& unzip protoc.zip -d /usr/local \
	&& chmod 755 /usr/local/bin/protoc \
	&& rm -rf "${TEMP}"

	

USER hypnos

RUN  echo \
		google.golang.org/protobuf/cmd/protoc-gen-go@v1.31 \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3 \
		entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc@master \
		github.com/yoheimuta/protolint/cmd/protolint@latest \
	| xargs -n1 go install

# Comment below lines if you don't want to use Zsh
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true \
	&& git clone https://github.com/zsh-users/zsh-autosuggestions         ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions     \
	&& git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting \
	&& git clone https://github.com/zsh-users/zsh-completions             ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions         \
	&& git clone https://github.com/supercrabtree/k                       ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/k
ADD .zshrc /home/hypnos
ENV SHELL /bin/zsh
