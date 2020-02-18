package kubernetes

func DeleteContainerByTaskID(taskID uint64) error {
	err1 := DeleteDeployByTaskID(taskID)
	err2 := DeleteServiceByTaskID(taskID)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
