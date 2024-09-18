#          _                   _            _             _           _            _
#         / /\                /\ \         /\ \     _    /\ \        /\ \         /\_\
#        / /  \              /  \ \       /  \ \   /\_\ /  \ \      /  \ \       / / /  _
#       / / /\ \            / /\ \ \     / /\ \ \_/ / // /\ \ \    / /\ \ \     / / /  /\_\
#      / / /\ \ \          / / /\ \ \   / / /\ \___/ // / /\ \_\  / / /\ \ \   / / /__/ / /
#     / / /  \ \ \        / / /  \ \_\ / / /  \/____// / /_/ / / / / /  \ \_\ / /\_____/ /
#    / / /___/ /\ \      / / /   / / // / /    / / // / /__\/ / / / /   / / // /\_______/
#   / / /_____/ /\ \    / / /   / / // / /    / / // / /_____/ / / /   / / // / /\ \ \
#  / /_________/\ \ \  / / /___/ / // / /    / / // / /\ \ \  / / /___/ / // / /  \ \ \
# / / /_       __\ \_\/ / /____\/ // / /    / / // / /  \ \ \/ / /____\/ // / /    \ \ \
# \_\___\     /____/_/\/_________/ \/_/     \/_/ \/_/    \_\/\/_________/ \/_/      \_\_\
# Developed by AonrokZa1392
# ติดต่อแก้สคริปได้ที่เฟส AonrokZa1392 ไม่เข้ารหัสไฟล์ support ตลอดการใช้งาน

# Builder stage
FROM golang:1.22 AS builder
WORKDIR /app

COPY . . 
RUN go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o ./inu-backyard ./cmd/http_server/main.go

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o ./auto_migration ./cmd/auto_migration/main.go

# Runner stage
FROM alpine:3.19 AS runner
WORKDIR /app

# Install python and necessary packages
RUN apk add --no-cache python3 py3-pip py3-scipy py3-scikit-learn

# Copy over the requirements file and install Python dependencies
COPY --from=builder /app/requirements.txt /app/
RUN pip3 install --no-deps -r /app/requirements.txt --break-system-packages

# Copy binaries and other necessary files from the builder stage
COPY --from=builder /app/inu-backyard /app/inu-backyard
COPY --from=builder /app/predict.py /app/predict.py
COPY --from=builder /app/config.yml /app/config.yml
COPY --from=builder /app/auto_migration /app/auto_migration

EXPOSE 3001

CMD [ "/app/auto_migration" , "&&" , "/app/inu-backyard" ]