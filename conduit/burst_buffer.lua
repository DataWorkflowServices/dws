function HelloWorld()
	print ("Hello World")
end

function GetGoodCR()
	local cr = [[
apiVersion: v1alpha1
kind: ClientMount
spec:
  desiredState: mounted
  node: node1
  mounts:
    - type: lustre
      device:
        type: lustre
      mountPath: /lus
      options: --mgs 127.0.0.1
      targetType: directory
]]
	return cr
end

function GetBadCR()
        local cr = [[
apiVersion: v1alpha1
kind: ClientMount
spec:
  desiredState: badValue
  node: node1
  mounts:
    - type: lustre
      device:
        type: lustre
      mountPath: /lus
      options: --mgs 127.0.0.1
      targetType: directory
]]
        return cr
end
