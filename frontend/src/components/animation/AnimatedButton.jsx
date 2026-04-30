import { motion } from "motion/react";

export default function AnimatedButton({
  children,
  stiffness = 300,
  damping = 15,
  ...props  
}) {
  return (
    <motion.button
      whileHover={{ scale: 1.05, y: -2 }}
      whileTap={{ scale: 0.9, y: 1 }}
      transition={{ type: "spring", stiffness, damping }}
      {...props} 
    >
      {children}
    </motion.button>
  );
}